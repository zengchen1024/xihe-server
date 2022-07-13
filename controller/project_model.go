package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type projectCreateModel struct {
	Name      string `json:"name" required:"true"`
	Desc      string `json:"desc" required:"true"`
	Type      string `json:"type" required:"true"`
	CoverId   string `json:"cover_id" required:"true"`
	Protocol  string `json:"protocol" required:"true"`
	Training  string `json:"training" required:"true"`
	Inference string `json:"inference" required:"true"`
}

func (p *projectCreateModel) toCmd() (cmd app.ProjectCreateCmd, err error) {
	cmd.Name, err = domain.NewProjName(p.Name)
	if err != nil {
		return
	}

	cmd.Type, err = domain.NewRepoType(p.Type)
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

	cmd.Training, err = domain.NewTrainingSDK(p.Training)
	if err != nil {
		return
	}

	cmd.Inference, err = domain.NewInferenceSDK(p.Inference)
	if err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type projectUpdateModel struct {
	Name    *string `json:"name"`
	Desc    *string `json:"desc"`
	Type    *string `json:"type"`
	CoverId *string `json:"cover_id"`
	// json [] will be converted to []string
	Tags []string `json:"tags"`
}

func (p *projectUpdateModel) toCmd() (cmd app.ProjectUpdateCmd, err error) {
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

	if p.Type != nil {
		cmd.Type, err = domain.NewRepoType(*p.Type)
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
