package domain

import (
	td "github.com/opensourceways/xihe-training-center/domain"
)

type TrainingParam = td.Training
type KeyValue = td.KeyValue
type Compute = td.Compute

type UserTraining struct {
	Owner Account

	Training
}

type Training struct {
	Id string

	TrainingParam

	// following fileds is not under the controlling of version
	Job Job
}

type Job struct {
	Endpoint string
	Id       string
}

type TrainingInfo struct {
	Status   string
	Duration int
}
