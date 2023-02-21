package repositoryimpl

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongoCmdSet = "$set"
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	ObjectIdFilter(s string) (bson.M, error)

	AppendElemMatchToFilter(array string, exists bool, cond, filter bson.M)

	GetDoc(
		ctx context.Context, filterOfDoc, project bson.M,
		result interface{},
	) error

	GetDocs(
		ctx context.Context, filterOfDoc, project bson.M,
		result interface{},
	) error

	NewDocIfNotExist(
		ctx context.Context, filterOfDoc, docInfo bson.M,
	) (string, error)

	PushArrayElem(
		ctx context.Context, array string,
		filterOfDoc, value bson.M,
	) error

	UpdateArrayElem(
		ctx context.Context, array string,
		filterOfDoc, filterOfArray, updateCmd bson.M,
		version int, t int64,
	) (bool, error)

	UpdateDoc(
		ctx context.Context, filterOfDoc, update bson.M, op string, version int,
	) error
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
