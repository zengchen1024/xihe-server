package authing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Authing/authing-go-sdk/lib/authentication"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
)

var (
	cli      *authentication.Client
	endpoint string
)

func Init(appId, secret, authEndpoint string) {
	cli = authentication.NewClient(appId, secret)

	if strings.HasSuffix(authEndpoint, "/") {
		authEndpoint = strings.TrimSuffix(authEndpoint, "/")
	}
	endpoint = authEndpoint
}

func NewAuthingUser() authing.User {
	return user{}
}

type user struct{}

func (impl user) GetByAccessToken(accessToken string) (userInfo authing.UserInfo, err error) {
	/*
		respStr, err := cli.GetUserInfoByAccessToken(accessToken)
		if err != nil {
			return
		}
	*/

	url := endpoint + "/oidc/me?access_token=" + accessToken
	resp, err := cli.SendHttpRequest(url, http.MethodGet, nil, nil)

	// TODO: delete
	fmt.Printf("%s\n", resp)

	var loginInfo struct {
		Name    string `json:"name,omitempty"`
		Picture string `json:"picture,omitempty"`
		Email   string `json:"email,omitempty"`
	}

	err = json.Unmarshal(resp, &loginInfo)
	if err != nil {
		return
	}

	if userInfo.Name, err = domain.NewAccount(loginInfo.Name); err != nil {
		return
	}

	if userInfo.Email, err = domain.NewEmail(loginInfo.Email); err != nil {
		return
	}

	if userInfo.AvatarId, err = domain.NewAvatarId(loginInfo.Picture); err != nil {
		return
	}

	//TODO

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

	err = json.Unmarshal([]byte(respStr), &v)
	if err != nil {
		return
	}

	at := v.AccessToken
	if at == "" {
		err = errors.New("no access token")

		return
	}

	info, err := impl.GetByAccessToken(at)
	if err != nil {
		return
	}

	login.IDToken = v.IdToken
	login.UserInfo = info
	login.AccessToken = at

	return
}
