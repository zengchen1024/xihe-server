package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
)

type userBasicInfoUpdateRequest struct {
	AvatarId string `json:"avatar_id"`
	Bio      string `json:"bio"`
}

func (req *userBasicInfoUpdateRequest) toCmd() (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	cmd.Bio, err = domain.NewBio(req.Bio)
	if err != nil {
		return
	}

	cmd.AvatarId, err = domain.NewAvatarId(req.AvatarId)

	return
}

type userCreateRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (req *userCreateRequest) toCmd(info authing.UserInfo) (cmd app.UserCreateCmd, err error) {
	cmd.Account = info.Name
	cmd.Email = info.Email
	cmd.Bio = info.Bio
	cmd.AvatarId = info.AvatarId

	cmd.Password, err = domain.NewPassword(req.Password)
	if err != nil {
		return
	}

	err = cmd.Validate()

	return
}
