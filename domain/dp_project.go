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
	reName         = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	config         = Config{}
	reResourceName = reName
	reEmail        = regexp.MustCompile("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z]{2,6}$")
)

func Init(cfg Config) {
	config = cfg
}

type Config struct {
	Resource ResourceConfig
	User     UserConfig
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

type UserConfig struct {
	MaxNicknameLength int
	MaxBioLength      int
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

// Name
type ResourceName interface {
	ResourceName() string
}

// ResourceName
type ProjName interface {
	ProjName() string

	ResourceName
}

func NewProjName(v string) (ProjName, error) {
	if err := checkResourceName(v, ResourceProject); err != nil {
		return nil, err
	}

	return projName(v), nil

}

type projName string

func (r projName) ProjName() string {
	return string(r)
}

func (r projName) ResourceName() string {
	return string(r)
}

type ModelName interface {
	ModelName() string

	ResourceName
}

func NewModelName(v string) (ModelName, error) {
	if err := checkResourceName(v, ResourceModel); err != nil {
		return nil, err
	}

	return modelName(v), nil
}

type modelName string

func (r modelName) ModelName() string {
	return string(r)
}

func (r modelName) ResourceName() string {
	return string(r)
}

type DatasetName interface {
	DatasetName() string

	ResourceName
}

func NewDatasetName(v string) (DatasetName, error) {
	if err := checkResourceName(v, ResourceDataset); err != nil {
		return nil, err
	}

	return datasetName(v), nil
}

type datasetName string

func (r datasetName) DatasetName() string {
	return string(r)
}

func (r datasetName) ResourceName() string {
	return string(r)
}

func checkResourceName(v, prefix string) error {
	max := config.Resource.MaxNameLength
	min := config.Resource.MinNameLength

	if n := len(v); n > max || n < min {
		return fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if strings.HasPrefix(strings.ToLower(v), prefix) {
		return fmt.Errorf("the name should not start with %s as prefix", prefix)
	}

	if !reResourceName.MatchString(v) {
		return errors.New("invalid name")
	}

	return nil
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
