package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *userCreateRequest) toCmd(accessToken string) (cmd app.UserCreateCmd, err error) {
	// TODO get user's account and email by access token

	cmd.Account, err = domain.NewAccount(req.Account)
	if err != nil {
		return
	}

	cmd.Email, err = domain.NewEmail(req.Email)
	if err != nil {
		return
	}

	cmd.Password, err = domain.NewPassword(req.Password)
	if err != nil {
		return
	}

	err = cmd.Validate()

	return
}
