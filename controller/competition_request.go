package controller

import "github.com/opensourceways/xihe-server/domain"

type competitionAddRelatedProjectRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"project_name"`
}

func (req *competitionAddRelatedProjectRequest) toInfo() (
	owner domain.Account, name domain.ResourceName, err error,
) {
	if owner, err = domain.NewAccount(req.Owner); err != nil {
		return
	}

	name, err = domain.NewResourceName(req.Name)

	return
}
