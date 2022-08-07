package app

import (
	"fmt"

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
	items, err := s.repo.List(
		cmd.Owner,
		repository.ProjectListOption{Name: cmd.From.Name},
	)
	if err != nil {
		return
	}

	str := cmd.From.Name.ProjName()
	if n := len(items); n > 0 {
		str = fmt.Sprintf("%s%d", str, n)
	}

	name, err := domain.NewProjName(str)
	if err != nil {
		return
	}

	v := cmd.toProject(name)
	p, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	pid, err := s.pr.Fork(cmd.From.Id, platform.RepoOption{
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

	// TODO: webhook
	s.toProjectDTO(&p, &dto)

	return
}
