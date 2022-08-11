package gitlab

import (
	"strconv"

	sdk "github.com/xanzy/go-gitlab"

	"github.com/opensourceways/xihe-server/domain/platform"
)

type UserInfo struct {
	Token     string
	Namespace string
}

func NewRepositoryService(v UserInfo) platform.Repository {
	return &repository{user: v}
}

type repository struct {
	user UserInfo
}

func (r *repository) New(repo platform.RepoOption) (string, error) {
	cli, err := sdk.NewClient(r.user.Token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return "", err
	}

	ns, err := strconv.Atoi(r.user.Namespace)
	if err != nil {
		return "", err
	}

	name := repo.Name.ResourceName()
	des := repo.Desc.ProjDesc()
	b := true

	v, _, err := cli.Projects.CreateProject(&sdk.CreateProjectOptions{
		Name:                 &name,
		Path:                 &name,
		NamespaceID:          &ns,
		Description:          &des,
		InitializeWithReadme: &b,
	})

	if err != nil {
		return "", err
	}

	return strconv.Itoa(v.ID), nil
}

func (r *repository) Fork(srcRepoId string, repo platform.RepoOption) (string, error) {
	cli, err := sdk.NewClient(r.user.Token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return "", err
	}

	repoId, err := strconv.Atoi(srcRepoId)
	if err != nil {
		return "", err
	}

	name := repo.Name.ResourceName()
	des := repo.Desc.ProjDesc()
	b := true

	v, _, err := cli.Projects.ForkProject(repoId, &sdk.ForkProjectOptions{
		Name:                          &name,
		Path:                          &name,
		Description:                   &des,
		MergeRequestDefaultTargetSelf: &b,
	})

	if err != nil {
		return "", err
	}

	return strconv.Itoa(v.ID), nil
}
