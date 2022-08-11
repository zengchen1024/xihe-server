package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	RepoTypePublic  = "public"
	RepoTypePrivate = "priviate"
)

var (
	reProjName = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
)

// RepoType
type RepoType interface {
	RepoType() string
}

func NewRepoType(v string) (RepoType, error) {
	if v != RepoTypePublic && v != RepoTypePrivate {
		return nil, errors.New("unknown repo type")
	}

	return repoType(v), nil
}

type repoType string

func (r repoType) RepoType() string {
	return string(r)
}

// ProjName
type ProjName interface {
	ProjName() string
}

func NewProjName(v string) (ProjName, error) {
	return newResourceName(v, "project")
}

func NewModelName(v string) (ProjName, error) {
	return newResourceName(v, "model")
}

func NewDatasetName(v string) (ProjName, error) {
	return newResourceName(v, "dataset")
}

func newResourceName(v, prefix string) (ProjName, error) {
	if n := len(v); n > 50 || n < 5 {
		return nil, errors.New("name's length should be between 5 to 30 ")
	}

	if strings.HasPrefix(strings.ToLower(v), prefix) {
		return nil, fmt.Errorf("the name should not start with %s as prefix", prefix)
	}

	if !reProjName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return projName(v), nil
}

type projName string

func (r projName) ProjName() string {
	return string(r)
}

// ProjDesc
type ProjDesc interface {
	ProjDesc() string
}

func NewProjDesc(v string) (ProjDesc, error) {
	if len(v) > 100 || v == "" {
		return nil, errors.New("the length of desc should be between 1 to 100")
	}

	return projDesc(v), nil
}

type projDesc string

func (r projDesc) ProjDesc() string {
	return string(r)
}

// TrainingPlatform
type CoverId interface {
	CoverId() string
}

func NewConverId(v string) (CoverId, error) {
	if v == "" {
		return nil, errors.New("invalid cover id")
	}

	return coverId(v), nil
}

type coverId string

func (c coverId) CoverId() string {
	return string(c)
}

// ProtocolName
type ProtocolName interface {
	ProtocolName() string
}

func NewProtocolName(v string) (ProtocolName, error) {
	// TODO: limited value

	return protocolName(v), nil
}

type protocolName string

func (r protocolName) ProtocolName() string {
	return string(r)
}

// ProjType
type ProjType interface {
	ProjType() string
}

func NewProjType(v string) (ProjType, error) {
	if v != "Gradio" && v != "Static" {
		return nil, errors.New("unsupported project type")
	}

	return projType(v), nil
}

type projType string

func (r projType) ProjType() string {
	return string(r)
}

// TrainingPlatform
type TrainingPlatform interface {
	TrainingPlatform() string
}

func NewTrainingPlatform(v string) (TrainingPlatform, error) {
	if v != "ModelArts" {
		return nil, errors.New("unsupport training platform")
	}

	return trainingPlatform(v), nil
}

type trainingPlatform string

func (r trainingPlatform) TrainingPlatform() string {
	return string(r)
}
