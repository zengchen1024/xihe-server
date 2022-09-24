package domain

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"
)

var config configuration

func Init(r *ResourceConfig, u *UserConfig) {
	config = configuration{
		*r, *u,
	}
}

type configuration struct {
	ResourceConfig
	UserConfig
}

type ResourceConfig struct {
	covers           sets.String
	protocols        sets.String
	projectType      sets.String
	trainingPlatform sets.String

	MaxNameLength         int `json:"max_name_length"`
	MinNameLength         int `json:"min_name_length"`
	MaxDescLength         int `json:"max_desc_length"`
	MaxRelatedResourceNum int `json:"max_related_resource_num"`

	Covers           []string `json:"covers"            required:"true"`
	Protocols        []string `json:"protocols"         required:"true"`
	ProjectType      []string `json:"project_type"      required:"true"`
	TrainingPlatform []string `json:"training_platform" required:"true"`
}

func (r *ResourceConfig) SetDefault() {
	if r.MaxNameLength <= 0 {
		r.MaxNameLength = 50
	}

	if r.MinNameLength <= 0 {
		r.MinNameLength = 5
	}

	if r.MaxDescLength <= 0 {
		r.MaxDescLength = 100
	}

	if r.MaxRelatedResourceNum <= 0 {
		r.MaxRelatedResourceNum = 5
	}
}

func (r *ResourceConfig) Validate() error {
	if r.MaxNameLength < (r.MinNameLength + 10) {
		return errors.New("invalid name length")
	}

	r.covers = sets.NewString(r.Covers...)
	r.protocols = sets.NewString(r.Protocols...)
	r.projectType = sets.NewString(r.ProjectType...)
	r.trainingPlatform = sets.NewString(r.TrainingPlatform...)

	return nil
}

func (r *ResourceConfig) hasCover(v string) bool {
	return r.covers.Has(v)
}

func (r *ResourceConfig) hasProtocol(v string) bool {
	return r.protocols.Has(v)
}

func (r *ResourceConfig) hasProjectType(v string) bool {
	return r.projectType.Has(v)
}

func (r *ResourceConfig) hasPlatform(v string) bool {
	return r.trainingPlatform.Has(v)
}

type UserConfig struct {
	MaxNicknameLength int `json:"max_nickname_length"`
	MaxBioLength      int `json:"max_bio_length"`
}

func (u *UserConfig) SetDefault() {
	if u.MaxNicknameLength == 0 {
		u.MaxNicknameLength = 20
	}

	if u.MaxBioLength == 0 {
		u.MaxBioLength = 200
	}
}
