package authing

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserInfo struct {
	Name     domain.Account
	Email    domain.Email
	Bio      domain.Bio
	AvatarId domain.AvatarId
}

type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

type User interface {
	GetByAccessToken(accessToken string) (UserInfo, error)
	GetByCode(code string) (Login, error)
}
