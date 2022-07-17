package domain

import "errors"

const (
	RepoTypePublic  = "public"
	RepoTypePrivate = "priviate"
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
	// TODO: limited length for name

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
	// TODO: limited length for name

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
	// TODO: limited value

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
	// TODO: limited value

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
	// TODO: limited value

	return trainingPlatform(v), nil
}

type trainingPlatform string

func (r trainingPlatform) TrainingPlatform() string {
	return string(r)
}
