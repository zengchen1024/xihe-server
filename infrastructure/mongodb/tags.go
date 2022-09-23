package mongodb

import (
	"context"
	"errors"
	"sort"

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

	if err := col.sortTags(&v); err != nil {
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

func (col tags) sortTags(v *dResourceTags) error {
	items := v.Items
	orders := v.Orders

	for i := range items {
		if _, ok := orders[items[i].Domain]; !ok {
			return errors.New("can't sort tags")
		}
	}

	sort.Slice(items, func(i, j int) bool {

		a := orders[items[i].Domain]
		b := orders[items[j].Domain]

		if a != b {
			return a < b
		}

		return items[i].Order < items[j].Order
	})

	return nil
}
