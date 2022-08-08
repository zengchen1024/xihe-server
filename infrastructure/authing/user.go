package authing

import (
	"encoding/json"
	"errors"

	"github.com/Authing/authing-go-sdk/lib/authentication"

	"github.com/opensourceways/xihe-server/domain"
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

func (impl user) GetByAccessToken(accessToken string) (userInfo authing.UserInfo, err error) {
	respStr, err := cli.GetUserInfoByAccessToken(accessToken)
	if err != nil {
		return
	}

	var loginInfo struct {
		Birthdate  string `json:"birthdate,omitempty"`
		Gender     string `json:"gender,omitempty"`
		Name       string `json:"name,omitempty"`
		Nickname   string `json:"nickname,omitempty"`
		UserName   string `json:"username,omitempty"`
		Picture    string `json:"picture,omitempty"`
		UpdatedAT  string `json:"updated_at,omitempty"`
		Website    string `json:"website,omitempty"`
		ExternalID string `json:"external_id,omitempty"`
		Sub        string `json:"sub,omitempty"`
		Email      string `json:"email,omitempty"`
		// EmailVerified       bool   `json:"email_verified,omitempty"`
		// PhoneNumber         string `json:"phone_number,omitempty"`
		// PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
	}

	err = json.Unmarshal([]byte(respStr), &loginInfo)
	if err != nil {
		return
	}

	if userInfo.Name, err = domain.NewAccount(loginInfo.UserName); err != nil {
		return
	}

	if userInfo.Email, err = domain.NewEmail(loginInfo.Email); err != nil {
		return
	}

	//TODO
	if userInfo.Bio, err = domain.NewBio(loginInfo.Email); err != nil {
		return
	}

	return
}

func (impl user) GetByCode(code string) (login authing.Login, err error) {
	respStr, err := cli.GetAccessTokenByCode(code)
	if err != nil {
		return
	}

	var accessToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		IdToken      string `json:"id_token"`
		Scope        string `json:"scope"`
		CreatedAt    int64  `json:"created_at"`
		AuthCode     string `json:"code"`
	}

	err = json.Unmarshal([]byte(respStr), &accessToken)
	if err != nil {
		return
	}

	if accessToken.AccessToken == "" {
		err = errors.New("no access token")

		return
	}

	info, err := impl.GetByAccessToken(accessToken.AccessToken)
	if err != nil {
		return
	}

	login.UserInfo = info
	login.AccessToken = accessToken.AccessToken

	return
}
