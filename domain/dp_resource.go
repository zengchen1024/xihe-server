package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	resourceProject = "project"
	resourceDataset = "dataset"
	resourceModel   = "model"

	SortTypeUpdateTime    = "update_time"
	SortTypeFirstLetter   = "first_letter"
	SortTypeDownloadCount = "download_count"
)

var (
	reName         = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reResourceName = reName

	ResourceTypeProject = resourceType(resourceProject)
	ResourceTypeModel   = resourceType(resourceModel)
	ResourceTypeDataset = resourceType(resourceDataset)
)

// Name
type ResourceName interface {
	ResourceName() string
	ResourceType() ResourceType
	FirstLetterOfName() byte
}

// ProjName
type ProjName interface {
	ProjName() string

	ResourceName
}

func GenProjName(v string) (ProjName, error) {
	name, err := genResourceName(v, ResourceTypeProject)
	if err != nil {
		return nil, err
	}

	return projName(name), nil
}

func NewProjName(v string) (ProjName, error) {
	if err := checkResourceName(v, ResourceTypeProject); err != nil {
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

func (r projName) FirstLetterOfName() byte {
	s := strings.TrimPrefix(string(r), ResourceTypeProject.PrefixToName())

	return s[0]
}

// ModelName
type ModelName interface {
	ModelName() string

	ResourceName
}

func GenModelName(v string) (ModelName, error) {
	name, err := genResourceName(v, ResourceTypeModel)
	if err != nil {
		return nil, err
	}

	return modelName(name), nil
}

func NewModelName(v string) (ModelName, error) {
	if err := checkResourceName(v, ResourceTypeModel); err != nil {
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

func (r modelName) FirstLetterOfName() byte {
	s := strings.TrimPrefix(string(r), ResourceTypeModel.PrefixToName())

	return s[0]
}

// DatasetName
type DatasetName interface {
	DatasetName() string

	ResourceName
}

func GenDatasetName(v string) (DatasetName, error) {
	name, err := genResourceName(v, ResourceTypeDataset)
	if err != nil {
		return nil, err
	}

	return datasetName(name), nil
}

func NewDatasetName(v string) (DatasetName, error) {
	if err := checkResourceName(v, ResourceTypeDataset); err != nil {
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

func (r datasetName) FirstLetterOfName() byte {
	s := strings.TrimPrefix(string(r), ResourceTypeDataset.PrefixToName())

	return s[0]
}

func genResourceName(v string, t ResourceType) (string, error) {
	prefix := t.PrefixToName()

	max := config.MaxNameLength - len(prefix)
	min := config.MinNameLength

	if n := len(v); n > max || n < min {
		return "", fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if strings.HasPrefix(strings.ToLower(v), t.ResourceType()) {
		return "", fmt.Errorf(
			"the name should not start with %s as prefix", t.ResourceType(),
		)
	}

	if !reResourceName.MatchString(v) {
		return "", errors.New("invalid name")
	}

	return prefix + v, nil
}

func checkResourceName(v string, t ResourceType) error {
	max := config.MaxNameLength
	min := config.MinNameLength

	if n := len(v); n > max || n < min {
		return fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !strings.HasPrefix(strings.ToLower(v), t.PrefixToName()) {
		return fmt.Errorf(
			"the name should start with %s as prefix", t.PrefixToName(),
		)
	}

	if !reResourceName.MatchString(v) {
		return errors.New("invalid name")
	}

	return nil
}

// ResourceType
type ResourceType interface {
	ResourceType() string
	PrefixToName() string
}

func NewResourceType(v string) (ResourceType, error) {
	b := v == resourceProject ||
		v == resourceModel ||
		v == resourceDataset
	if b {
		return resourceType(v), nil
	}

	return nil, errors.New("invalid resource type")
}

type resourceType string

func (r resourceType) ResourceType() string {
	return string(r)
}

func (r resourceType) PrefixToName() string {
	return string(r) + "-"
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

// SortType
type SortType interface {
	SortType() string
}

func NewSortType(v string) (SortType, error) {
	b := v != SortTypeUpdateTime &&
		v != SortTypeFirstLetter &&
		v != SortTypeDownloadCount

	if b {
		return nil, errors.New("invliad sort type")
	}

	return sortType(v), nil
}

type sortType string

func (s sortType) SortType() string {
	return string(s)
}
