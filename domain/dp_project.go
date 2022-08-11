package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	RepoTypePublic  = "public"
	RepoTypePrivate = "priviate"

	ResourceProject = "project"
	ResourceDataset = "dataset"
	ResourceModel   = "model"
)

var (
	reProjName = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	config     = Config{}
)

func Init(cfg Config) {
	config = cfg
}

type Config struct {
	Resource ResourceConfig
}

type ResourceConfig struct {
	MaxNameLength int
	MinNameLength int
	MaxDescLength int

	Covers           sets.String
	Protocols        sets.String
	ProjectType      sets.String
	TrainingPlatform sets.String
}

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
	return newResourceName(v, ResourceProject)
}

func NewModelName(v string) (ProjName, error) {
	return newResourceName(v, ResourceModel)
}

func NewDatasetName(v string) (ProjName, error) {
	return newResourceName(v, ResourceDataset)
}

func newResourceName(v, prefix string) (ProjName, error) {
	max := config.Resource.MaxNameLength
	min := config.Resource.MinNameLength

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
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
	max := config.Resource.MaxDescLength
	if len(v) > max || v == "" {
		return nil, fmt.Errorf("the length of desc should be between 1 to %d", max)
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
	if !config.Resource.Covers.Has(v) {
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
	if !config.Resource.Protocols.Has(v) {
		return nil, errors.New("unsupported protocol")
	}

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
	if !config.Resource.ProjectType.Has(v) {
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
	if !config.Resource.TrainingPlatform.Has(v) {
		return nil, errors.New("unsupport training platform")
	}

	return trainingPlatform(v), nil
}

type trainingPlatform string

func (r trainingPlatform) TrainingPlatform() string {
	return string(r)
}
