package authing

import (
	"github.com/Authing/authing-go-sdk/lib/authentication"

	"github.com/opensourceways/xihe-server/domain/authing"
)

var cli *authentication.Client

func Init(appId, secret string) {
	cli = authentication.NewClient(appId, secret)
}

func NewAuthingUser() authing.User {
	return user{}
}

type user struct{}

func (impl user) GetByAccessToken(accessToken string) (userinfo authing.UserInfo, err error) {
	return
}

func (impl user) GetByCode(code string) (login authing.Login, err error) {
	return
}
