package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewTagsMapper(name string) repositories.TagsMapper {
	return tags{name}
}

type tags struct {
	collectionName string
}

func (col tags) List(domainNames []string) ([]repositories.DomainTagsDo, error) {
	var v []dResourceTags

	if err := col.listTags(domainNames, &v); err != nil || len(v) == 0 {
		return nil, err
	}

	orders := map[string]int{}
	for i, n := range domainNames {
		orders[n] = i
	}

	items := v[0].Items
	r := make([]repositories.DomainTagsDo, len(items))
	for i := range items {
		r[i] = col.toDomainTagsDO(&items[i])
	}

	return r, nil
}

func (col tags) listTags(domainNames []string, v *[]dResourceTags) error {
	fieldItemsRef := "$" + fieldItems

	project := bson.M{
		fieldItems: bson.M{"$filter": bson.M{
			"input": fieldItemsRef,
			"cond": func() bson.M {
				conds := bson.A{}

				if len(domainNames) > 0 {
					conds = append(conds, inCondForArrayElem(
						fieldName, domainNames,
					))
				}

				return condForArrayElem(conds)
			}(),
		}},
	}

	pipeline := bson.A{
		bson.M{"$project": project},
	}

	return withContext(func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, v)
	})
}

func (col tags) toDomainTagsDO(doc *dDomainTags) (do repositories.DomainTagsDo) {
	do.Name = doc.Name
	do.Domain = doc.Domain

	tags := doc.Tags
	do.Items = make([]repositories.TagsDo, len(tags))
	for i := range tags {
		do.Items[i] = repositories.TagsDo{
			Kind:  tags[i].Kind,
			Items: tags[i].Tags,
		}
	}

	return
}
