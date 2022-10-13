package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type trainingCreateResp struct {
	Id string `json:"id"`
}

type trainingLogResp struct {
	LogURL string `json:"log_url"`
}

type TrainingCreateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`

	CodeDir  string `json:"code_dir"`
	BootFile string `json:"boot_file"`

	Hypeparameters []KeyValue      `json:"hyperparameter"`
	Env            []KeyValue      `json:"evn"`
	Models         []TrainingInput `json:"models"`
	Datasets       []TrainingInput `json:"datasets"`

	Compute Compute `json:"compute"`
}

func (req *TrainingCreateRequest) toCmd(cmd *app.TrainingCreateCmd) (err error) {
	if cmd.Name, err = domain.NewTrainingName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewTrainingDesc(req.Desc); err != nil {
		return
	}

	if cmd.CodeDir, err = domain.NewDirectory(req.CodeDir); err != nil {
		return
	}

	if cmd.BootFile, err = domain.NewFilePath(req.BootFile); err != nil {
		return
	}

	if cmd.Compute, err = req.Compute.toCompute(); err != nil {
		return
	}

	if cmd.Env, err = req.toKeyValue(req.Env); err != nil {
		return
	}

	cmd.Hypeparameters, err = req.toKeyValue(req.Hypeparameters)

	return
}

func (req *TrainingCreateRequest) toKeyValue(kv []KeyValue) (r []domain.KeyValue, err error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r = make([]domain.KeyValue, n)
	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}

type Compute struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

func (c *Compute) toCompute() (r domain.Compute, err error) {
	if c.Type == "" || c.Version == "" || c.Flavor == "" {
		err = errors.New("invalid compute info")

		return
	}

	if r.Type, err = domain.NewComputeType(c.Type); err != nil {
		return
	}

	if r.Version, err = domain.NewComputeVersion(c.Version); err != nil {
		return
	}

	if r.Flavor, err = domain.NewComputeFlavor(c.Flavor); err != nil {
		return
	}

	return
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *KeyValue) toKeyValue() (r domain.KeyValue, err error) {
	if kv.Key == "" {
		err = errors.New("invalid key value")

		return
	}

	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

type TrainingInput struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
	File  string `json:"File"`
}

func (t *TrainingInput) toModelInput() (r domain.Input, name domain.ModelName, err error) {
	if err = t.toInput(&r); err != nil {
		return
	}

	if name, err = domain.NewModelName(t.Name); err != nil {
		return
	}

	r.Value.Type = domain.ResourceTypeModel

	return
}

func (t *TrainingInput) toDatasetInput() (r domain.Input, name domain.DatasetName, err error) {
	if err = t.toInput(&r); err != nil {
		return
	}

	if name, err = domain.NewDatasetName(t.Name); err != nil {
		return
	}

	r.Value.Type = domain.ResourceTypeDataset

	return
}

func (t *TrainingInput) toInput(r *domain.Input) (err error) {
	if r.Key, err = domain.NewCustomizedKey(t.Key); err != nil {
		return
	}

	if r.Value.User, err = domain.NewAccount(t.Owner); err != nil {
		return
	}

	r.Value.File = t.File

	return
}

func (ctl *TrainingController) setProjectInfo(
	ctx *gin.Context, cmd *app.TrainingCreateCmd,
	user domain.Account, projectId string,
) (ok bool) {
	v, err := ctl.project.GetSummary(user, projectId)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	name, ok := v.Name.(domain.ProjName)
	if !ok {
		ctl.sendRespWithInternalError(ctx, newResponseError(
			errors.New("it is not a project name"),
		))

		return
	}

	cmd.User = user
	cmd.ProjectId = projectId
	cmd.ProjectName = name
	cmd.ProjectRepoId = v.RepoId
	ok = true

	return
}

func (ctl *TrainingController) setModelsInput(
	ctx *gin.Context, cmd *app.TrainingCreateCmd,
	inputs []TrainingInput,
) (ok bool) {
	p := make([]repository.ModelSummaryListOption, 0, len(inputs))
	m := sets.NewString()
	tinputs := make([]domain.Input, len(inputs))

	index := func(u domain.Account, n string) string {
		return u.Account() + n
	}

	for i := range inputs {
		v := &inputs[i]

		ti, name, err := v.toModelInput()
		if err != nil {
			if err != nil {
				ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

				return
			}
		}
		tinputs[i] = ti

		s := index(ti.Value.User, v.Name)
		if m.Has(s) {
			continue
		}
		m.Insert(s)

		p = append(p, repository.ModelSummaryListOption{
			Owner: ti.Value.User,
			Name:  name,
		})

	}

	v, err := ctl.model.ListSummary(p)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	if len(v) != len(p) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("some models does not exist"),
			))

			return
		}
	}

	ids := map[string]string{}
	for i := range v {
		item := &v[i]

		ids[index(item.Owner, item.Name.ResourceName())] = item.RepoId
	}

	for i := range inputs {
		v := &tinputs[i].Value

		v.RepoId = ids[index(v.User, inputs[i].Name)]
	}

	cmd.Inputs = append(cmd.Inputs, tinputs...)
	ok = true

	return
}

func (ctl *TrainingController) setDatasetsInput(
	ctx *gin.Context, cmd *app.TrainingCreateCmd,
	inputs []TrainingInput,
) (ok bool) {
	p := make([]repository.DatasetSummaryListOption, 0, len(inputs))
	m := sets.NewString()
	tinputs := make([]domain.Input, len(inputs))

	index := func(u domain.Account, n string) string {
		return u.Account() + n
	}

	for i := range inputs {
		v := &inputs[i]

		ti, name, err := v.toDatasetInput()
		if err != nil {
			if err != nil {
				ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

				return
			}
		}
		tinputs[i] = ti

		s := index(ti.Value.User, v.Name)
		if m.Has(s) {
			continue
		}
		m.Insert(s)

		p = append(p, repository.DatasetSummaryListOption{
			Owner: ti.Value.User,
			Name:  name,
		})

	}

	v, err := ctl.dataset.ListSummary(p)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	if len(v) != len(p) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("some datasets does not exist"),
			))

			return
		}
	}

	ids := map[string]string{}
	for i := range v {
		item := &v[i]

		ids[index(item.Owner, item.Name.ResourceName())] = item.RepoId
	}

	for i := range inputs {
		v := &tinputs[i].Value

		v.RepoId = ids[index(v.User, inputs[i].Name)]
	}

	cmd.Inputs = append(cmd.Inputs, tinputs...)
	ok = true

	return
}
