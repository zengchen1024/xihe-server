package gitlab

import (
	"strconv"

	sdk "github.com/xanzy/go-gitlab"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

var admin *administrator

func NewUserSerivce() platform.User {
	return admin
}

func Init(endpoint, token string) error {
	v, err := sdk.NewClient(token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return err
	}

	admin = &administrator{v}

	return nil
}

type administrator struct {
	cli *sdk.Client
}

func (m *administrator) New(u platform.UserOption) (r domain.PlatformUser, err error) {
	name := u.Name.Account()
	email := u.Email.Email()
	pass := u.Password.Password()
	b := true

	v, _, err := m.cli.Users.CreateUser(&sdk.CreateUserOptions{
		Name:             &name,
		Email:            &email,
		Username:         &name,
		Password:         &pass,
		SkipConfirmation: &b,
	})

	if err != nil {
		return
	}

	r.Id = strconv.Itoa(v.ID)
	r.NamespaceId = strconv.Itoa(v.NamespaceID)

	return
}

func (m *administrator) NewToken(u domain.PlatformUser) (string, error) {
	uid, err := strconv.Atoi(u.Id)
	if err != nil {
		return "", err
	}

	name := "___"
	scope := []string{"api"}

	v, _, err := m.cli.Users.CreatePersonalAccessToken(
		uid, &sdk.CreatePersonalAccessTokenOptions{
			Name:   &name,
			Scopes: &scope,
		},
	)

	if err != nil {
		return "", err
	}

	return v.Token, nil
}
