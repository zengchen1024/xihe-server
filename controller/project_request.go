package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type projectCreateRequest struct {
	Owner    string `json:"owner" required:"true"`
	Name     string `json:"name" required:"true"`
	Desc     string `json:"desc"`
	Type     string `json:"type" required:"true"`
	CoverId  string `json:"cover_id" required:"true"`
	Protocol string `json:"protocol" required:"true"`
	Training string `json:"training" required:"true"`
	RepoType string `json:"repo_type" required:"true"`
}

func (p *projectCreateRequest) toCmd() (cmd app.ProjectCreateCmd, err error) {
	if cmd.Owner, err = domain.NewAccount(p.Owner); err != nil {
		return
	}

	if cmd.Name, err = domain.GenProjName(p.Name); err != nil {
		return
	}

	if cmd.Type, err = domain.NewProjType(p.Type); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewResourceDesc(p.Desc); err != nil {
		return
	}

	if cmd.CoverId, err = domain.NewConverId(p.CoverId); err != nil {
		return
	}

	if cmd.Protocol, err = domain.NewProtocolName(p.Protocol); err != nil {
		return
	}

	if cmd.RepoType, err = domain.NewRepoType(p.RepoType); err != nil {
		return
	}

	if cmd.Training, err = domain.NewTrainingPlatform(p.Training); err != nil {
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
}

func (p *projectUpdateRequest) toCmd() (cmd app.ProjectUpdateCmd, err error) {
	if p.Name != nil {
		if cmd.Name, err = domain.NewProjName(*p.Name); err != nil {
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = domain.NewResourceDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.RepoType != nil {
		if cmd.RepoType, err = domain.NewRepoType(*p.RepoType); err != nil {
			return
		}
	}

	if p.CoverId != nil {
		if cmd.CoverId, err = domain.NewConverId(*p.CoverId); err != nil {
			return
		}
	}

	return
}

type projectDetail struct {
	Liked bool `json:"liked"`

	*app.ProjectDTO
}
