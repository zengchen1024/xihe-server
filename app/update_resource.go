package app

import "k8s.io/apimachinery/pkg/util/sets"

type ResourceTagsUpdateCmd struct {
	ToAdd    []string
	ToRemove []string
}

func (cmd *ResourceTagsUpdateCmd) toTags(old []string) ([]string, bool) {
	tags := sets.NewString(old...)

	if len(cmd.ToAdd) > 0 {
		tags.Insert(cmd.ToAdd...)
	}

	if len(cmd.ToRemove) > 0 {
		tags.Delete(cmd.ToRemove...)
	}

	if tags.Equal(sets.NewString(old...)) {
		return nil, false
	}

	return tags.UnsortedList(), true
}
