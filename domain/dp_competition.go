package domain

import (
	"errors"
	"net/url"
)

const (
	competitionTeamRoleLeader = "leader"

	competitionPhaseFinal       = "final"
	competitionPhasePreliminary = "preliminary"

	competitionStatusDone       = "done"
	competitionStatusPreparing  = "preparing"
	competitionStatusInProgress = "in-progress"

	competitionIdentityStudent   = "student"
	competitionIdentityTeacher   = "teacher"
	competitionIdentityDeveloper = "developer"
)

// CompetitionPhase
type CompetitionPhase interface {
	CompetitionPhase() string
}

func NewCompetitionPhase(v string) (CompetitionPhase, error) {
	if v == competitionPhasePreliminary || v == competitionPhaseFinal {
		return competitionPhase(v), nil
	}

	return nil, errors.New("invalid competition phase")
}

type competitionPhase string

func (r competitionPhase) CompetitionPhase() string {
	return string(r)
}

// CompetitionStatus
type CompetitionStatus interface {
	CompetitionStatus() string
}

func NewCompetitionStatus(v string) (CompetitionStatus, error) {
	b := v == competitionStatusPreparing ||
		v == competitionStatusInProgress ||
		v == competitionStatusDone

	if b {
		return competitionStatus(v), nil
	}

	return nil, errors.New("invalid competition status")
}

type competitionStatus string

func (r competitionStatus) CompetitionStatus() string {
	return string(r)
}

// CompetitionName
type CompetitionName interface {
	CompetitionName() string
}

func NewCompetitionName(v string) (CompetitionName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionName(v), nil
}

type competitionName string

func (r competitionName) CompetitionName() string {
	return string(r)
}

// CompetitionDesc
type CompetitionDesc interface {
	CompetitionDesc() string
}

func NewCompetitionDesc(v string) (CompetitionDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionDesc(v), nil
}

type competitionDesc string

func (r competitionDesc) CompetitionDesc() string {
	return string(r)
}

// CompetitionDuration
type CompetitionDuration interface {
	CompetitionDuration() string
}

func NewCompetitionDuration(v string) (CompetitionDuration, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionDuration(v), nil
}

type competitionDuration string

func (r competitionDuration) CompetitionDuration() string {
	return string(r)
}

// CompetitionBonus
type CompetitionBonus interface {
	CompetitionBonus() int
}

func NewCompetitionBonus(v int) (CompetitionBonus, error) {
	if v == 0 {
		return nil, errors.New("empty value")
	}

	return competitionBonus(v), nil
}

type competitionBonus int

func (r competitionBonus) CompetitionBonus() int {
	return int(r)
}

//
type CompetitionHost interface {
	CompetitionHost() string
}

func NewCompetitionHost(v string) (CompetitionHost, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitionHost(v), nil
}

type competitionHost string

func (r competitionHost) CompetitionHost() string {
	return string(r)
}

// URL
type URL interface {
	URL() string
}

func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpURL(v), nil
}

type dpURL string

func (r dpURL) URL() string {
	return string(r)
}

// Phone
type Phone interface {
	Phone() string
}

func NewPhone(v string) (Phone, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return phone(v), nil
}

type phone string

func (r phone) Phone() string {
	return string(r)
}

// CompetitionIdentity
type CompetitionIdentity interface {
	CompetitionIdentity() string
}

func NewcompetitionIdentity(v string) (CompetitionIdentity, error) {
	b := v == competitionIdentityStudent ||
		v == competitionIdentityTeacher ||
		v == competitionIdentityDeveloper ||
		v == ""

	if !b {
		return nil, errors.New("invalid competition identity")
	}

	return competitionIdentity(v), nil
}

type competitionIdentity string

func (r competitionIdentity) CompetitionIdentity() string {
	return string(r)
}

// Province
type Province interface {
	Province() string
}

func NewProvince(v string) (Province, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return province(v), nil
}

type province string

func (r province) Province() string {
	return string(r)
}

// City
type City interface {
	City() string
}

func NewCity(v string) (City, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return city(v), nil
}

type city string

func (r city) City() string {
	return string(r)
}

// CompetitorName
type CompetitorName interface {
	CompetitorName() string
}

func NewCompetitorName(v string) (CompetitorName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return competitorName(v), nil
}

type competitorName string

func (r competitorName) CompetitorName() string {
	return string(r)
}

// TeamName
type TeamName interface {
	TeamName() string
}

func NewTeamName(v string) (TeamName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return teamName(v), nil
}

type teamName string

func (r teamName) TeamName() string {
	return string(r)
}

// TeamRole
type TeamRole interface {
	TeamRole() string
}

func NewTeamRole(v string) (TeamRole, error) {
	if v == "" || v == competitionTeamRoleLeader {
		return teamRole(v), nil
	}

	return nil, errors.New("invalid team role")
}

type teamRole string

func (r teamRole) TeamRole() string {
	return string(r)
}
