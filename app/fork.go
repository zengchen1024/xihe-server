package app

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectForkCmd struct {
	From  domain.Project
	Owner domain.Account
}

func (cmd *ProjectForkCmd) toProject(name domain.ProjName) domain.Project {
	p := &cmd.From
	return domain.Project{
		Owner:    cmd.Owner,
		Name:     name,
		Desc:     p.Desc,
		Type:     p.Type,
		CoverId:  p.CoverId,
		RepoType: p.RepoType,
		Protocol: p.Protocol,
		Training: p.Training,
		Tags:     p.Tags,
	}
}

func (s projectService) Fork(cmd *ProjectForkCmd) (dto ProjectDTO, err error) {
	name, err := s.genForkedRepoName(cmd.Owner, cmd.From.Name)
	if err != nil {
		return
	}

	v := cmd.toProject(name)

	p, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	pid, err := s.pr.Fork(cmd.From.RepoId, platform.RepoOption{
		Name: name,
		Desc: cmd.From.Desc,
	})
	if err != nil {
		return
	}

	p.RepoId = pid

	p, err = s.repo.Save(&p)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	return
}

func (s projectService) genForkedRepoName(
	owner domain.Account, srcName domain.ProjName,
) (domain.ProjName, error) {
	items, err := s.repo.List(
		owner,
		repository.ProjectListOption{Name: srcName},
	)
	if err != nil {
		return nil, err
	}

	names := sets.NewString()
	for i := range items {
		names.Insert(items[i].Name.ProjName())
	}

	str := srcName.ProjName()
	n := len(items)
	for {
		if !names.Has(str) {
			break
		}

		str = fmt.Sprintf("%s%d", str, n)
		n += 1
	}

	return domain.NewProjName(str)
}
