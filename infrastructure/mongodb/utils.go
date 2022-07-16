package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errDocExists    = errors.New("doc exists")
	errDocNotExists = errors.New("doc doesn't exist")
)

type dbError struct {
	error
}

func isDBError(err error) bool {
	_, b := err.(dbError)

	return b
}

func toUID(docId interface{}) (string, error) {
	if v, ok := docId.(primitive.ObjectID); ok {
		return v.Hex(), nil
	}

	return "", errors.New("retrieve id failed")
}

func newId() string {
	return primitive.NewObjectID().Hex()
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	err = json.Unmarshal(v, &m)

	return
}

func appendElemMatchToFilter(array string, exists bool, cond, filter bson.M) {
	match := bson.M{"$elemMatch": cond}

	if exists {
		filter[array] = match
	} else {
		filter[array] = bson.M{"$not": match}
	}
}

func (cli *client) newDocIfNotExist(
	ctx context.Context, collection string,
	filterOfDoc, docInfo bson.M,
) (string, error) {
	upsert := true
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{"$setOnInsert": docInfo},
		&options.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", dbError{err}
	}

	if r.UpsertedID == nil {
		return "", errDocExists
	}

	v, _ := toUID(r.UpsertedID)

	return v, nil
}

func (cli *client) pushArrayElem(
	ctx context.Context,
	collection, array string,
	filterOfDoc, value bson.M,
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{"$push": bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) getArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray bson.M,
	project bson.M, result interface{},
) error {
	ma := map[string]bson.M{}
	if len(filterOfArray) > 0 {
		ma[array] = filterOfArray
	}

	return cli.getArraysElem(ctx, collection, filterOfDoc, ma, project, result)
}

func (cli *client) getArraysElem(
	ctx context.Context, collection string,
	filterOfDoc bson.M, filterOfArrays map[string]bson.M,
	project bson.M, result interface{},
) error {
	m := map[string]func() bson.M{}
	for k, v := range filterOfArrays {
		m[k] = func() bson.M {
			return conditionTofilterArray(v)
		}
	}

	return cli.getArraysElemsByCustomizedCond(
		ctx, collection, filterOfDoc, m,
		project, result,
	)
}

func (cli *client) getArraysElemsByCustomizedCond(
	ctx context.Context, collection string,
	filterOfDoc bson.M, filterOfArrays map[string]func() bson.M,
	project bson.M, result interface{},
) error {
	pipeline := bson.A{bson.M{"$match": filterOfDoc}}

	if len(filterOfArrays) > 0 {
		project1 := bson.M{}

		for array, cond := range filterOfArrays {
			project1[array] = bson.M{"$filter": bson.M{
				"input": fmt.Sprintf("$%s", array),
				"cond":  cond(),
			}}
		}

		for k, v := range project {
			s := k
			if i := strings.Index(k, "."); i >= 0 {
				s = k[:i]
			}
			if _, ok := filterOfArrays[s]; !ok {
				project1[k] = v
			}
		}

		pipeline = append(pipeline, bson.M{"$project": project1})
	}

	if len(project) > 0 {
		pipeline = append(pipeline, bson.M{"$project": project})
	}

	col := cli.collection(collection)
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return cursor.All(ctx, result)
}

func conditionTofilterArray(filterOfArray bson.M) bson.M {
	cond := make(bson.A, 0, len(filterOfArray))
	for k, v := range filterOfArray {
		cond = append(cond, bson.M{"$eq": bson.A{"$$this." + k, v}})
	}

	if len(filterOfArray) == 1 {
		return cond[0].(bson.M)
	}

	return bson.M{"$and": cond}
}
