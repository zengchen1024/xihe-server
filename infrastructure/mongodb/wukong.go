package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewWuKongMapper(name string) repositories.WuKongMapper {
	return wukong{
		collectionName: name,
	}
}

type wukong struct {
	collectionName string
}

func (col wukong) ListSamples(docId string, nums []int) ([]string, error) {
	project := bson.M{
		fieldSamples: bson.M{"$filter": bson.M{
			"input": "$" + fieldSamples,
			"cond": func() bson.M {
				return inCondForArrayElem(fieldNum, nums)
			}(),
		}},
	}

	pipeline := bson.A{
		bson.M{"$match": resourceIdFilter(docId)},
		bson.M{"$project": project},
	}

	var v []dWuKong

	err := withContext(func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})
	if err != nil || len(v) == 0 {
		return nil, err
	}

	doc := v[0]

	r := make([]string, len(doc.Samples))
	for i := range doc.Samples {
		r[i] = doc.Samples[i].Name
	}

	return r, nil
}

func (col wukong) ListPictures(docId string, opt *repositories.WuKongPictureListOptionDO) (
	do repositories.WuKongPicturesDO, err error,
) {
	p := 0
	if opt.PageNum > 1 {
		p = opt.CountPerPage * (opt.PageNum - 1)
	}
	fieldRef := "$" + fieldPictures

	project := bson.M{
		fieldPictures: bson.M{"$slice": bson.A{
			fieldRef, p, opt.CountPerPage,
		}},
		"total": bson.M{
			"$cond": bson.M{
				"if":   bson.M{"$isArray": fieldRef},
				"then": bson.M{"$size": fieldRef},
				"else": 0,
			},
		},
	}

	pipeline := bson.A{
		bson.M{"$match": resourceIdFilter(docId)},
		bson.M{"$project": project},
	}

	var v []struct {
		Total    int           `bson:"total"`
		Pictures []pictureInfo `bson:"pictures"`
	}

	err = withContext(func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})
	if err != nil || len(v) == 0 {
		return
	}

	doc := v[0]
	do.Total = doc.Total

	do.Pictures = make([]repositories.WuKongPictureInfoDO, len(doc.Pictures))
	for i := range doc.Pictures {
		col.toPictureDO(&doc.Pictures[i], &do.Pictures[i])
	}

	return
}

func (col wukong) toPictureDO(p *pictureInfo, do *repositories.WuKongPictureInfoDO) {
	do.Link = p.Link
	do.Desc = p.Desc
	do.Style = p.Style
}
