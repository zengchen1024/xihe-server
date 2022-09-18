package platform

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserOption struct {
	Name     domain.Account
	Email    domain.Email
	Password domain.Password
}

type User interface {
	New(UserOption) (domain.PlatformUser, error)
	NewToken(domain.PlatformUser) (string, error)
}

type RepoOption struct {
	Name     domain.ResourceName
	RepoType domain.RepoType
}

func (r *RepoOption) IsNotEmpty() bool {
	return r.Name != nil || r.RepoType != nil
}

type Repository interface {
	New(repo *RepoOption) (string, error)
	Fork(srcRepoId string, Name domain.ResourceName) (string, error)
	Update(repoId string, repo *RepoOption) error
}
