package domain

import (
	"errors"
)

const (
	activityTypeLike   = "like"
	activityTypeCreate = "create"
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

func NewActivityTypeLike() ActivityType {
	return activityType(activityTypeLike)
}

func NewActivityTypeCreate() ActivityType {
	return activityType(activityTypeCreate)
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

	ResourceObj
}
