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
	cmd.Owner = req.Owner

	cmd.Name, err = domain.NewProjName(req.Name)
	if err != nil {
		return
	}

	cmd.Desc, err = domain.NewProjDesc(req.Desc)
	if err != nil {
		return
	}

	cmd.Protocol, err = domain.NewProtocolName(req.Protocol)
	if err != nil {
		return
	}

	cmd.RepoType, err = domain.NewRepoType(req.RepoType)
	if err != nil {
		return
	}

	err = cmd.Validate()

	return
}
