package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func isDocNotExists(err error) bool {
	return errors.Is(err, errDocNotExists)
}

func isDocExists(err error) bool {
	return errors.Is(err, errDocExists)
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

func objectIdFilter(s string) (bson.M, error) {
	v, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id": v,
	}, nil
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	err = json.Unmarshal(v, &m)

	return
}

func isErrNoDocuments(err error) bool {
	return err.Error() == mongo.ErrNoDocuments.Error()
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

func (cli *client) replaceDoc(
	ctx context.Context, collection string,
	filterOfDoc, docInfo bson.M,
) (string, error) {
	upsert := true

	r, err := cli.collection(collection).ReplaceOne(
		ctx, filterOfDoc, docInfo,
		&options.ReplaceOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", dbError{err}
	}

	if r.UpsertedID == nil {
		return "", nil
	}

	v, _ := toUID(r.UpsertedID)
	return v, nil
}

func (cli *client) updateDoc(
	ctx context.Context, collection string,
	filterOfDoc, update bson.M, version int,
) error {
	filterOfDoc[fieldVersion] = version
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			"$set": update,
			"$inc": bson.M{fieldVersion: 1},
		},
	)

	if err != nil {
		return err
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) getDoc(
	ctx context.Context, collection string,
	filterOfDoc, project bson.M, result interface{},
) error {
	var sr *mongo.SingleResult
	col := cli.collection(collection)
	if len(project) > 0 {
		sr = col.FindOne(ctx, filterOfDoc, &options.FindOneOptions{
			Projection: project,
		})
	} else {
		sr = col.FindOne(ctx, filterOfDoc)
	}

	if err := sr.Decode(result); err != nil {
		if isErrNoDocuments(err) {
			return errDocNotExists
		}

		return err
	}

	return nil
}

func (cli *client) getDocs(
	ctx context.Context, collection string,
	filterOfDoc, project bson.M, result interface{},
) error {
	col := cli.collection(collection)

	var cursor *mongo.Cursor
	var err error
	if len(project) > 0 {
		cursor, err = col.Find(ctx, filterOfDoc, &options.FindOptions{
			Projection: project,
		})
	} else {
		cursor, err = col.Find(ctx, filterOfDoc)
	}

	if err != nil {
		return err
	}

	return cursor.All(ctx, result)
}

func (cli *client) addToSimpleArray(
	ctx context.Context, collection, array string,
	filterOfDoc, value interface{},
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{"$addToSet": bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	if r.ModifiedCount == 0 {
		return errDocExists
	}

	return nil
}

func (cli *client) removeFromSimpleArray(
	ctx context.Context,
	collection, array string,
	filterOfDoc, value interface{},
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{"$pull": bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
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

func (cli *client) updateArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, updateCmd bson.M, version int,
) (bool, error) {
	cmd := bson.M{}
	for k, v := range updateCmd {
		cmd[fmt.Sprintf("%s.$[i].%s", array, k)] = v
	}

	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}
	arrayFilter["i."+fieldVersion] = version

	col := cli.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			"$set": cmd,
			"$inc": bson.M{fmt.Sprintf("%s.$[i].%s", array, fieldVersion): 1},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if r.MatchedCount == 0 {
		return false, errDocNotExists
	}

	return r.ModifiedCount > 0, nil
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

func condFieldOfArrayElem(key string) string {
	return "$$this." + key
}

func eqCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{"$eq": bson.A{condFieldOfArrayElem(key), value}}
}

func matchCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{
		"$regexMatch": bson.M{
			"input": condFieldOfArrayElem(key),
			"regex": value,
		},
	}
}

func condForArrayElem(conds bson.A) bson.M {
	n := len(conds)
	if n > 1 {
		return bson.M{"$and": conds}
	}

	if n == 1 {
		return conds[0].(bson.M)
	}

	return bson.M{
		"$toBool": 1,
	}
}

func conditionTofilterArray(filterOfArray bson.M) bson.M {
	cond := make(bson.A, 0, len(filterOfArray))
	for k, v := range filterOfArray {
		cond = append(cond, eqCondForArrayElem(k, v))
	}

	return condForArrayElem(cond)
}
