package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForProjectController(
	rg *gin.RouterGroup,
	repo repository.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	ctl := ProjectController{
		repo:    repo,
		model:   model,
		dataset: dataset,
		s:       app.NewProjectService(repo, activity, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/project", ctl.Create)
	rg.PUT("/v1/project/:owner/:id", ctl.Update)
	rg.GET("/v1/project/:owner/:name", ctl.Get)
	rg.GET("/v1/project/:owner", ctl.List)

	rg.POST("/v1/project/:owner/:id", ctl.Fork)

	rg.PUT("/v1/project/:owner/:id/relation", ctl.AddRelatedResource)
	rg.DELETE("/v1/project/:owner/:id/relation", ctl.RemoveRelatedResource)
}

type ProjectController struct {
	baseController

	repo repository.Project
	s    app.ProjectService

	model   repository.Model
	dataset repository.Dataset

	newPlatformRepository func(string, string) platform.Repository
}

// @Summary Create
// @Description create project
// @Tags  Project
// @Param	body	body 	projectCreateRequest	true	"body of creating project"
// @Accept json
// @Produce json
// @Router /v1/project [post]
func (ctl *ProjectController) Create(ctx *gin.Context) {
	req := projectCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create project for other user",
		))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Create(&cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary Update
// @Description update project
// @Tags  Project
// @Param	owner	path	string			true	"owner of project"
// @Param	id	path	string			true	"id of project"
// @Param	body	body 	projectUpdateRequest	true	"body of updating project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{id} [put]
func (ctl *ProjectController) Update(ctx *gin.Context) {
	req := projectUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't update project for other user",
		))

		return
	}

	proj, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Update(&proj, &cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary Get
// @Description get project
// @Tags  Project
// @Param	owner	path	string	true	"owner of project"
// @Param	name	path	string	true	"name of project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{name} [get]
func (ctl *ProjectController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewProjName(ctx.Param("name"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	proj, err := ctl.s.GetByName(owner, name)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	if (visitor || pl.isNotMe(owner)) && proj.RepoType != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(proj))
	}
}

// @Summary List
// @Description list project
// @Tags  Project
// @Param	owner	path	string			true	"owner of project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner} [get]
func (ctl *ProjectController) List(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd, err := ctl.getListParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if visitor || pl.isNotMe(owner) {
		if cmd.RepoType == nil {
			cmd.RepoType, _ = domain.NewRepoType(domain.RepoTypePublic)
		} else {
			if cmd.RepoType.RepoType() != domain.RepoTypePublic {
				ctx.JSON(http.StatusOK, newResponseData(nil))

				return
			}
		}
	}

	projs, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(projs))
}

func (ctl *ProjectController) getListParameter(
	ctx *gin.Context,
) (cmd app.ResourceListCmd, err error) {
	if v := ctl.getQueryParameter(ctx, "name"); v != "" {
		cmd.Name = v
	}

	if v := ctl.getQueryParameter(ctx, "repo_type"); v != "" {
		if cmd.RepoType, err = domain.NewRepoType(v); err != nil {
			return
		}
	}

	return
}

// @Summary Fork
// @Description fork project
// @Tags  Project
// @Param	owner	path	string	true	"owner of forked project"
// @Param	id	path	string	true	"id of forked project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{id} [post]
func (ctl *ProjectController) Fork(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if !pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "no need to fork project of yourself",
		))

		return
	}

	proj, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	// TODO maybe the private project can be forked by special user.

	if proj.RepoType.RepoType() != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	data, err := ctl.s.Fork(&app.ProjectForkCmd{
		From:  proj,
		Owner: pl.DomainAccount(),
	}, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

// @Summary AddRelatedResource
// @Description add related resource to project
// @Tags  Project
// @Param	owner	path	string				true	"owner of project"
// @Param	id	path	string				true	"id of project"
// @Param	body	body 	relatedResourceModifyRequest	true	"body of updating project"
// @Accept json
// @Success 200 {object} app.ResourceDTO
// @Router /v1/project/{owner}/{id}/relation [put]
func (ctl *ProjectController) AddRelatedResource(ctx *gin.Context) {
	ctl.relatedResource(ctx, true)
}

func (ctl *ProjectController) addRelatedResource(
	ctx *gin.Context, proj *domain.Project, cmd *domain.ResourceObj,
) {
	var f func(*domain.Project, *domain.ResourceIndex) error
	var data interface{}
	var err error

	switch cmd.ResourceType.ResourceType() {
	case domain.ResourceModel:
		data, err = ctl.model.Get(cmd.ResourceOwner, cmd.ResourceId)
		f = ctl.s.AddRelatedModel

	case domain.ResourceDataset:
		data, err = ctl.dataset.Get(cmd.ResourceOwner, cmd.ResourceId)
		f = ctl.s.AddRelatedDataset

	default:
		err = errors.New("unsupported related resource")
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	index := domain.ResourceIndex{
		ResourceOwner: cmd.ResourceOwner,
		ResourceId:    cmd.ResourceId,
	}
	if err = f(proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(convertToRelatedResource(data)))
}

// @Summary RemoveRelatedResource
// @Description remove related resource from project
// @Tags  Project
// @Param	owner	path	string				true	"owner of project"
// @Param	id	path	string				true	"id of project"
// @Param	body	body 	relatedResourceModifyRequest	true	"body of updating project"
// @Accept json
// @Success 204
// @Router /v1/project/{owner}/{id}/relation [delete]
func (ctl *ProjectController) RemoveRelatedResource(ctx *gin.Context) {
	ctl.relatedResource(ctx, false)
}

func (ctl *ProjectController) removeRelatedResource(
	ctx *gin.Context, proj *domain.Project, cmd *domain.ResourceObj,
) {
	var err error
	var f func(*domain.Project, *domain.ResourceIndex) error

	switch cmd.ResourceType.ResourceType() {
	case domain.ResourceModel:
		f = ctl.s.RemoveRelatedModel

	case domain.ResourceDataset:
		f = ctl.s.RemoveRelatedDataset

	default:
		err = errors.New("unsupported related resource")
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	index := domain.ResourceIndex{
		ResourceOwner: cmd.ResourceOwner,
		ResourceId:    cmd.ResourceId,
	}
	if err = f(proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

func (ctl *ProjectController) relatedResource(ctx *gin.Context, add bool) {
	req := relatedResourceModifyRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't update project for other user",
		))

		return
	}

	proj, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	if add {
		ctl.addRelatedResource(ctx, &proj, &cmd)
	} else {
		ctl.removeRelatedResource(ctx, &proj, &cmd)
	}
}
