package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type datasetCreateRequest struct {
	Owner    string `json:"owner" required:"true"`
	Name     string `json:"name" required:"true"`
	Desc     string `json:"desc"`
	Protocol string `json:"protocol" required:"true"`
	RepoType string `json:"repo_type" required:"true"`
}

func (req *datasetCreateRequest) toCmd() (cmd app.DatasetCreateCmd, err error) {
	if cmd.Owner, err = domain.NewAccount(req.Owner); err != nil {
		return
	}

	if cmd.Name, err = domain.GenDatasetName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewProjDesc(req.Desc); err != nil {
		return
	}

	if cmd.Protocol, err = domain.NewProtocolName(req.Protocol); err != nil {
		return
	}

	if cmd.RepoType, err = domain.NewRepoType(req.RepoType); err != nil {
		return
	}

	err = cmd.Validate()

	return
}
