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
	return resourceType(ResourceProject)
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
	return resourceType(ResourceModel)
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
	return resourceType(ResourceDataset)
}

func genResourceName(v, prefix string) (string, error) {
	max := config.Resource.MaxNameLength - len(prefix)
	min := config.Resource.MinNameLength

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
	max := config.Resource.MaxNameLength
	min := config.Resource.MinNameLength

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

func ResourceTypeByName(n string) (ResourceType, error) {
	if strings.HasPrefix(n, ResourceProject) {
		return NewResourceType(ResourceProject)
	}

	if strings.HasPrefix(n, ResourceDataset) {
		return NewResourceType(ResourceDataset)
	}

	if strings.HasPrefix(n, ResourceModel) {
		return NewResourceType(ResourceModel)
	}

	return nil, errors.New("unknow resource")
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

// ResourceObj
type ResourceObj struct {
	ResourceOwner Account
	ResourceType  ResourceType
	ResourceId    string
}
