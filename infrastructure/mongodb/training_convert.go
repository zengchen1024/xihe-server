package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col training) toTrainingDoc(do *repositories.UserTrainingDO) (bson.M, error) {
	tdo := &do.TrainingDO
	c := &tdo.Compute

	docObj := trainingItem{
		Id:             do.Id,
		Name:           tdo.Name,
		Desc:           tdo.Desc,
		CodeDir:        tdo.CodeDir,
		BootFile:       tdo.BootFile,
		CreatedAt:      do.CreatedAt,
		Inputs:         col.toInputDoc(tdo.Inputs),
		Env:            col.toKeyValueDoc(tdo.Env),
		Hypeparameters: col.toKeyValueDoc(tdo.Hypeparameters),
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

func (impl training) toInputDoc(v []repositories.InputDO) []dInput {
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
