package domain

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_-/]+$")
	reFilePath  = regexp.MustCompile("^[a-zA-Z0-9_-/.]+$")
)

type TrainingConfig struct {
	MaxNameLength int
	MinNameLength int
	MaxDescLength int
}

// TrainingName
type TrainingName interface {
	TrainingName() string
}

func NewTrainingName(v string) (TrainingName, error) {
	max := config.Training.MaxNameLength
	min := config.Training.MinNameLength

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return trainingName(v), nil
}

type trainingName string

func (r trainingName) TrainingName() string {
	return string(r)
}

// TrainingDesc
type TrainingDesc interface {
	TrainingDesc() string
}

func NewTrainingDesc(v string) (TrainingDesc, error) {
	max := config.Training.MaxDescLength
	if len(v) > max {
		return nil, fmt.Errorf("the length of desc should be less than %d", max)
	}

	return trainingDesc(v), nil
}

type trainingDesc string

func (r trainingDesc) TrainingDesc() string {
	return string(r)
}

// Directory
type Directory interface {
	Directory() string
}

func NewDirectory(v string) (Directory, error) {
	if !reDirectory.MatchString(v) {
		return nil, errors.New("invalid directory")
	}

	return directory(v), nil
}

type directory string

func (r directory) Directory() string {
	return string(r)
}

// FilePath
type FilePath interface {
	FilePath() string
}

func NewFilePath(v string) (FilePath, error) {
	if !reFilePath.MatchString(v) {
		return nil, errors.New("invalid filePath")
	}

	return filePath(v), nil
}

type filePath string

func (r filePath) FilePath() string {
	return string(r)
}

// ComputeType
type ComputeType interface {
	ComputeType() string
}

func NewComputeType(v string) (ComputeType, error) {
	if v == "" {
		return nil, errors.New("invalid compute type")
	}

	return computeType(v), nil
}

type computeType string

func (r computeType) ComputeType() string {
	return string(r)
}

// ComputeVersion
type ComputeVersion interface {
	ComputeVersion() string
}

func NewComputeVersion(v string) (ComputeVersion, error) {
	if v == "" {
		return nil, errors.New("invalid compute version")
	}

	return computeVersion(v), nil
}

type computeVersion string

func (r computeVersion) ComputeVersion() string {
	return string(r)
}

// ComputeFlavor
type ComputeFlavor interface {
	ComputeFlavor() string
}

func NewComputeFlavor(v string) (ComputeFlavor, error) {
	if v == "" {
		return nil, errors.New("invalid compute flavor")
	}

	return computeFlavor(v), nil
}

type computeFlavor string

func (r computeFlavor) ComputeFlavor() string {
	return string(r)
}

// CustomizedKey
type CustomizedKey interface {
	CustomizedKey() string
}

func NewCustomizedKey(v string) (CustomizedKey, error) {
	if v == "" {
		return nil, errors.New("invalid key")
	}

	return customizedKey(v), nil
}

type customizedKey string

func (r customizedKey) CustomizedKey() string {
	return string(r)
}

// CustomizedValue
type CustomizedValue interface {
	CustomizedValue() string
}

func NewCustomizedValue(v string) (CustomizedValue, error) {
	if v == "" {
		return nil, errors.New("invalid key")
	}

	return customizedValue(v), nil
}

type customizedValue string

func (r customizedValue) CustomizedValue() string {
	return string(r)
}

// TrainingRegion
type TrainingRegion interface {
	TrainingRegion() string
}

func NewTrainingRegion(v string) (TrainingRegion, error) {
	if v == "" {
		return nil, errors.New("invalid key")
	}

	return trainingRegion(v), nil
}

type trainingRegion string

func (r trainingRegion) TrainingRegion() string {
	return string(r)
}
