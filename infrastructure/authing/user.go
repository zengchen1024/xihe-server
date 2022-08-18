package authing

import (
	"encoding/json"
	"errors"

	"github.com/Authing/authing-go-sdk/lib/authentication"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
)

var (
	cli *authentication.Client
)

func Init(appId, secret, redirectURI string) {
	cli = authentication.NewClient(appId, secret)
	cli.RedirectUri = redirectURI
}

func NewAuthingUser() authing.User {
	return user{}
}

type user struct{}

func (impl user) GetByAccessToken(accessToken string) (userInfo authing.UserInfo, err error) {
	if accessToken == "" {
		err = errors.New("no access token")

		return
	}

	resp, err := cli.GetUserInfoByAccessToken(accessToken)
	if err != nil {
		return
	}

	var v struct {
		Name    string `json:"username,omitempty"`
		Picture string `json:"picture,omitempty"`
		Email   string `json:"email,omitempty"`
	}

	if err = json.Unmarshal([]byte(resp), &v); err != nil {
		return
	}

	if userInfo.Name, err = domain.NewAccount(v.Name); err != nil {
		return
	}

	if userInfo.Email, err = domain.NewEmail(v.Email); err != nil {
		return
	}

	if userInfo.AvatarId, err = domain.NewAvatarId(v.Picture); err != nil {
		return
	}

	return
}

func (impl user) GetByCode(code string) (login authing.Login, err error) {
	respStr, err := cli.GetAccessTokenByCode(code)
	if err != nil {
		return
	}

	var v struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}

	if err = json.Unmarshal([]byte(respStr), &v); err != nil {
		return
	}

	if v.IdToken == "" {
		err = errors.New("no id token")

		return
	}

	info, err := impl.GetByAccessToken(v.AccessToken)
	if err == nil {
		login.IDToken = v.IdToken
		login.UserInfo = info
	}

	return
}
