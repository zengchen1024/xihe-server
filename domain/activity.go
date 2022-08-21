package domain

import (
	"errors"
)

const (
	ActivityTypeLike   = "like"
	ActivityTypeCreate = "create"
)

// ActivityType
type ActivityType interface {
	ActivityType() string
}

func NewActivityType(v string) (ActivityType, error) {
	if v != ActivityTypeLike && v != ActivityTypeCreate {
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

	ResourceObj
}
