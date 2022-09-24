package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	ResourceProject = "project"
	ResourceDataset = "dataset"
	ResourceModel   = "model"
)

var (
	reName         = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reResourceName = reName

	ResourceTypeProject = resourceType(ResourceProject)
	ResourceTypeModel   = resourceType(ResourceModel)
	ResourceTypeDataset = resourceType(ResourceDataset)
)

// Name
type ResourceName interface {
	ResourceName() string
	ResourceType() ResourceType
}

// ProjName
type ProjName interface {
	ProjName() string

	ResourceName
}

func GenProjName(v string) (ProjName, error) {
	name, err := genResourceName(v, ResourceProject)
	if err != nil {
		return nil, err
	}

	return projName(name), nil
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

func (r projName) ResourceType() ResourceType {
	return ResourceTypeProject
}

// ModelName
type ModelName interface {
	ModelName() string

	ResourceName
}

func GenModelName(v string) (ModelName, error) {
	name, err := genResourceName(v, ResourceModel)
	if err != nil {
		return nil, err
	}

	return modelName(name), nil
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

func (r modelName) ResourceType() ResourceType {
	return ResourceTypeModel
}

// DatasetName
type DatasetName interface {
	DatasetName() string

	ResourceName
}

func GenDatasetName(v string) (DatasetName, error) {
	name, err := genResourceName(v, ResourceDataset)
	if err != nil {
		return nil, err
	}

	return datasetName(name), nil
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

func (r datasetName) ResourceType() ResourceType {
	return ResourceTypeDataset
}

func genResourceName(v, prefix string) (string, error) {
	max := config.MaxNameLength - len(prefix)
	min := config.MinNameLength

	if n := len(v); n > max || n < min {
		return "", fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if strings.HasPrefix(strings.ToLower(v), prefix) {
		return "", fmt.Errorf("the name should not start with %s as prefix", prefix)
	}

	if !reResourceName.MatchString(v) {
		return "", errors.New("invalid name")
	}

	return prefix + "-" + v, nil
}

func checkResourceName(v, prefix string) error {
	max := config.MaxNameLength
	min := config.MinNameLength

	if n := len(v); n > max || n < min {
		return fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if prefix += "-"; !strings.HasPrefix(strings.ToLower(v), prefix) {
		return fmt.Errorf("the name should start with %s as prefix", prefix)
	}

	if !reResourceName.MatchString(v) {
		return errors.New("invalid name")
	}

	return nil
}

// ResourceType
type ResourceType interface {
	ResourceType() string
}

func NewResourceType(v string) (ResourceType, error) {
	b := v == ResourceProject ||
		v == ResourceModel ||
		v == ResourceDataset
	if b {
		return resourceType(v), nil
	}

	return nil, errors.New("invalid resource type")
}

type resourceType string

func (r resourceType) ResourceType() string {
	return string(r)
}

// ResourceDesc
type ResourceDesc interface {
	ResourceDesc() string
}

func NewResourceDesc(v string) (ResourceDesc, error) {
	max := config.MaxDescLength
	if len(v) > max || v == "" {
		return nil, fmt.Errorf("the length of desc should be between 1 to %d", max)
	}

	return resourceDesc(v), nil
}

type resourceDesc string

func (r resourceDesc) ResourceDesc() string {
	return string(r)
}
