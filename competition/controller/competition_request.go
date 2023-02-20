package controller

import (
	"github.com/opensourceways/xihe-server/competition/app"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type competitionAddRelatedProjectRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"project_name"`
}

func (req *competitionAddRelatedProjectRequest) toInfo() (
	owner types.Account, name types.ResourceName, err error,
) {
	if owner, err = types.NewAccount(req.Owner); err != nil {
		return
	}

	name, err = types.NewResourceName(req.Name)

	return
}

type competitorApplyRequest struct {
	Name     string            `json:"name"`
	City     string            `json:"city"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Identity string            `json:"identity"`
	Province string            `json:"province"`
	Detail   map[string]string `json:"detail"`
}

func (req *competitorApplyRequest) toCmd(user types.Account) (cmd app.CompetitorApplyCmd, err error) {
	if cmd.Name, err = domain.NewCompetitorName(req.Name); err != nil {
		return
	}

	if cmd.City, err = domain.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = types.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = domain.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = domain.NewcompetitionIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = domain.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}
