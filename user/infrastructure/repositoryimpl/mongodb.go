package repositoryimpl

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongoCmdSet  = "$set"
	mongoCmdPush = "$push"
)

type mongodbClient interface {
	IsDocExists(error) bool
	IsDocNotExists(error) bool

	ObjectIdFilter(s string) (bson.M, error)
	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)
	UpdateDoc(ctx context.Context, filterOfDoc, update bson.M, op string, version int) error
	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	err = json.Unmarshal(v, &m)

	return
}
