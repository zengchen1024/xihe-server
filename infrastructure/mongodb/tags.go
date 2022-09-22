package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func tagsDocFilter(s string) bson.M {
	return bson.M{
		fieldRType: s,
	}
}

func NewTagsMapper(name string) repositories.TagsMapper {
	return tags{name}
}

type tags struct {
	collectionName string
}

func (col tags) List(resourceType string) ([]repositories.DomainTagsDo, error) {
	var v dResourceTags

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			tagsDocFilter(resourceType), nil, &v,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return nil, err
	}

	items := v.Items
	dt := make(map[string][]repositories.TagsDo)
	for i := range items {
		item := &items[i]

		domain := item.Domain
		tags := repositories.TagsDo{
			Kind:  item.Kind,
			Items: item.Tags,
		}

		if dos, ok := dt[domain]; !ok {
			dt[domain] = []repositories.TagsDo{tags}
		} else {
			dos = append(dos, tags)
			dt[domain] = dos
		}
	}

	n := len(dt)
	r := make([]repositories.DomainTagsDo, n)
	i := 0
	for domain, tags := range dt {
		r[i] = repositories.DomainTagsDo{
			Domain: domain,
			Items:  tags,
		}

		i++
	}

	return r, nil
}
