package mongodb

import (
	"context"
	"encoding/json"
	"errors"

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
