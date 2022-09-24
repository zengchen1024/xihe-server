package domain

import (
	"errors"
)

const (
	activityTypeLike   = "like"
	activityTypeCreate = "create"
)

var (
	ActivityTypeLike   = activityType(activityTypeLike)
	ActivityTypeCreate = activityType(activityTypeCreate)
)

// ActivityType
type ActivityType interface {
	ActivityType() string
}

func NewActivityType(v string) (ActivityType, error) {
	if v != activityTypeLike && v != activityTypeCreate {
		return nil, errors.New("unknown activity type")
	}

	return activityType(v), nil
}

type activityType string

func (r activityType) ActivityType() string {
	return string(r)
}

// UserActivity
type UserActivity struct {
	Owner Account

	Activity
}

type Activity struct {
	Type ActivityType
	Time int64

	ResourceObject
}
