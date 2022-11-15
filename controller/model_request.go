package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type modelCreateRequest struct {
	Owner    string `json:"owner" required:"true"`
	Name     string `json:"name" required:"true"`
	Desc     string `json:"desc"`
	Protocol string `json:"protocol" required:"true"`
	RepoType string `json:"repo_type" required:"true"`
}

func (req *modelCreateRequest) toCmd() (cmd app.ModelCreateCmd, err error) {
	if cmd.Owner, err = domain.NewAccount(req.Owner); err != nil {
		return
	}

	if cmd.Name, err = domain.NewResourceName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewResourceDesc(req.Desc); err != nil {
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

type modelUpdateRequest struct {
	Name     *string `json:"name"`
	Desc     *string `json:"desc"`
	RepoType *string `json:"type"`
}

func (p *modelUpdateRequest) toCmd() (cmd app.ModelUpdateCmd, err error) {
	if p.Name != nil {
		if cmd.Name, err = domain.NewResourceName(*p.Name); err != nil {
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

	return
}

type modelDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*app.ModelDetailDTO
}

type modelsInfo struct {
	Owner    string `json:"owner"`
	AvatarId string `json:"avatar_id"`

	*app.ModelsDTO
}
