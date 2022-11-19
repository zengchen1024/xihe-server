package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col training) toTrainingDoc(do *repositories.UserTrainingDO) (bson.M, error) {
	cfg := &do.TrainingConfigDO
	c := &cfg.Compute

	docObj := trainingItem{
		Id:             do.Id,
		Name:           cfg.Name,
		Desc:           cfg.Desc,
		CodeDir:        cfg.CodeDir,
		BootFile:       cfg.BootFile,
		CreatedAt:      do.CreatedAt,
		Inputs:         col.toInputDoc(cfg.Inputs),
		EnableAim:      do.EnableAim,
		EnableOutput:   do.EnableOutput,
		Env:            col.toKeyValueDoc(cfg.Env),
		Hypeparameters: col.toKeyValueDoc(cfg.Hypeparameters),
		Compute: dCompute{
			Type:    c.Type,
			Flavor:  c.Flavor,
			Version: c.Version,
		},
	}
	return genDoc(docObj)
}

func (col training) toKeyValueDoc(kv []repositories.KeyValueDO) []dKeyValue {
	n := len(kv)
	if n == 0 {
		return nil
	}

	r := make([]dKeyValue, n)

	for i := range kv {
		r[i].Key = kv[i].Key
		r[i].Value = kv[i].Value
	}

	return r
}

func (col training) toInputDoc(v []repositories.InputDO) []dInput {
	n := len(v)
	if n == 0 {
		return nil
	}

	r := make([]dInput, n)

	for i := range v {
		item := &v[i]

		r[i] = dInput{
			Key:    item.Key,
			User:   item.User,
			Type:   item.Type,
			File:   item.File,
			RepoId: item.RepoId,
		}
	}

	return r
}

func (col training) toTrainingDetailDO(doc *dTraining) repositories.TrainingDetailDO {
	item := &doc.Items[0]

	return repositories.TrainingDetailDO{
		CreatedAt:        item.CreatedAt,
		Job:              col.toTrainingJobInfoDO(&item.Job),
		JobDetail:        col.toTrainingJobDetailDO(&item.JobDetail),
		TrainingConfigDO: col.toTrainingConfigDO(doc),
	}
}

func (col training) toTrainingConfigDO(doc *dTraining) repositories.TrainingConfigDO {
	item := &doc.Items[0]
	c := &item.Compute

	return repositories.TrainingConfigDO{
		ProjectName:    doc.ProjectName,
		ProjectRepoId:  doc.ProjectRepoId,
		Name:           item.Name,
		Desc:           item.Desc,
		CodeDir:        item.CodeDir,
		BootFile:       item.BootFile,
		Inputs:         col.toInputs(item.Inputs),
		EnableAim:      item.EnableAim,
		EnableOutput:   item.EnableOutput,
		Env:            col.toKeyValues(item.Env),
		Hypeparameters: col.toKeyValues(item.Hypeparameters),
		Compute: repositories.ComputeDO{
			Type:    c.Type,
			Flavor:  c.Flavor,
			Version: c.Version,
		},
	}
}

func (col training) toKeyValues(kv []dKeyValue) []repositories.KeyValueDO {
	n := len(kv)
	if n == 0 {
		return nil
	}

	r := make([]repositories.KeyValueDO, n)

	for i := range kv {
		r[i].Key = kv[i].Key
		r[i].Value = kv[i].Value
	}

	return r
}

func (col training) toInputs(v []dInput) []repositories.InputDO {
	n := len(v)
	if n == 0 {
		return nil
	}

	r := make([]repositories.InputDO, n)

	for i := range v {
		item := &v[i]

		r[i] = repositories.InputDO{
			Key:    item.Key,
			User:   item.User,
			Type:   item.Type,
			File:   item.File,
			RepoId: item.RepoId,
		}
	}

	return r
}

func (col training) toTrainingJobInfoDO(doc *dJobInfo) repositories.TrainingJobInfoDO {
	return repositories.TrainingJobInfoDO{
		Endpoint:  doc.Endpoint,
		JobId:     doc.JobId,
		LogDir:    doc.LogDir,
		AimDir:    doc.AimDir,
		OutputDir: doc.OutputDir,
	}
}

func (col training) toTrainingJobDetailDO(doc *dJobDetail) repositories.TrainingJobDetailDO {
	return repositories.TrainingJobDetailDO{
		Error:      doc.Error,
		Status:     doc.Status,
		Duration:   doc.Duration,
		LogPath:    doc.LogPath,
		AimPath:    doc.AimPath,
		OutputPath: doc.OutputPath,
	}
}
