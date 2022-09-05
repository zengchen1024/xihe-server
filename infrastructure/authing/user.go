package authing

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
)

var (
	cli *authentication.Client
)

func Init(appId, secret string) {
	cli = authentication.NewClient(appId, secret)
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

	var v struct {
		Name    string `json:"username,omitempty"`
		Picture string `json:"picture,omitempty"`
		Email   string `json:"email,omitempty"`
	}

	if err = getUserInfoByAccessToken(accessToken, &v); err != nil {
		return
	}

	fmt.Printf("user info = %v\n", v)
	
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

func (impl user) GetByCode(code, redirectURI string) (login authing.Login, err error) {
	var v struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}

	if err = getAccessTokenByCode(code, redirectURI, &v); err != nil {
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

func getAccessTokenByCode(code, redirectURI string, result interface{}) error {
	body := map[string]string{
		"client_id":     cli.AppId,
		"client_secret": cli.Secret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
	}

	value := make(url.Values)
	for k, v := range body {
		value.Add(k, v)
	}

	req, err := http.NewRequest(
		http.MethodPost, cli.Host+"/oidc/token",
		strings.NewReader(strings.TrimSpace(value.Encode())),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return sendHttpRequest(req, result)
}

func getUserInfoByAccessToken(accessToken string, result interface{}) error {
	url := cli.Host + "/oidc/me?access_token=" + accessToken

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return sendHttpRequest(req, result)
}

func sendHttpRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "xihe-server-authing")

	hc := utils.HttpClient{MaxRetries: 3}

	return hc.ForwardTo(req, result)
}
