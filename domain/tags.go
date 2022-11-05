package domain

import "k8s.io/apimachinery/pkg/util/sets"

type DomainTags struct {
	Name   string
	Domain string `json:"domain"`
	Items  []Tags `json:"items"`
}

func (t *DomainTags) GetKindsOfTags(tags []string) []string {
	r := make([]string, 0, len(t.Items))

	for i := range t.Items {
		if k := t.Items[i].getKindIfIncludes(tags); k != "" {
			r = append(r, k)
		}
	}

	return r
}

type Tags struct {
	Kind  string   `json:"kind"`
	Items []string `json:"items"`
}

func (t *Tags) getKindIfIncludes(tags []string) string {
	if t.Kind == "" {
		return ""
	}

	if len(tags) <= len(t.Items) {
		if sets.NewString(tags...).HasAny(t.Items...) {
			return t.Kind
		}

		return ""
	}

	if sets.NewString(t.Items...).HasAny(tags...) {
		return t.Kind
	}

	return ""
}
