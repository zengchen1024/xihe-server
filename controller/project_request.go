package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type projectCreateRequest struct {
	Name     string `json:"name" required:"true"`
	Desc     string `json:"desc"`
	Type     string `json:"type" required:"true"`
	CoverId  string `json:"cover_id" required:"true"`
	Protocol string `json:"protocol" required:"true"`
	Training string `json:"training" required:"true"`
	RepoType string `json:"repo_type" required:"true"`
}

func (p *projectCreateRequest) toCmd(owner string) (cmd app.ProjectCreateCmd, err error) {
	cmd.Owner = owner

	cmd.Name, err = domain.NewProjName(p.Name)
	if err != nil {
		return
	}

	cmd.Type, err = domain.NewProjType(p.Type)
	if err != nil {
		return
	}

	cmd.Desc, err = domain.NewProjDesc(p.Desc)
	if err != nil {
		return
	}

	cmd.CoverId, err = domain.NewConverId(p.CoverId)
	if err != nil {
		return
	}

	cmd.Protocol, err = domain.NewProtocolName(p.Protocol)
	if err != nil {
		return
	}

	cmd.RepoType, err = domain.NewRepoType(p.RepoType)
	if err != nil {
		return
	}

	cmd.Training, err = domain.NewTrainingPlatform(p.Training)
	if err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type projectUpdateRequest struct {
	Name     *string `json:"name"`
	Desc     *string `json:"desc"`
	RepoType *string `json:"type"`
	CoverId  *string `json:"cover_id"`
	// json [] will be converted to []string
	Tags []string `json:"tags"`
}

func (p *projectUpdateRequest) toCmd() (cmd app.ProjectUpdateCmd, err error) {
	if p.Name != nil {
		cmd.Name, err = domain.NewProjName(*p.Name)
		if err != nil {
			return
		}
	}

	if p.Desc != nil {
		cmd.Desc, err = domain.NewProjDesc(*p.Desc)
		if err != nil {
			return
		}
	}

	if p.RepoType != nil {
		cmd.RepoType, err = domain.NewRepoType(*p.RepoType)
		if err != nil {
			return
		}
	}

	if p.CoverId != nil {
		cmd.CoverId, err = domain.NewConverId(*p.CoverId)
		if err != nil {
			return
		}
	}

	if p.Tags != nil {
		// TODO check tags
	}

	return
}
