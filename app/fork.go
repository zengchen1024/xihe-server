package app

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectForkCmd struct {
	From      domain.Project
	Owner     domain.Account
	ValidTags []domain.DomainTags
}

func (cmd *ProjectForkCmd) toProject(name domain.ProjName, r *domain.Project) {
	p := &cmd.From
	*r = domain.Project{
		Owner:     cmd.Owner,
		Type:      p.Type,
		Protocol:  p.Protocol,
		Training:  p.Training,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		ProjectModifiableProperty: domain.ProjectModifiableProperty{
			Name:     p.Name,
			Desc:     p.Desc,
			CoverId:  p.CoverId,
			RepoType: p.RepoType,
			Tags:     p.Tags,
		},
	}

	h := ResourceTagsUpdateCmd{
		All: cmd.ValidTags,
	}

	r.TagKinds = h.genTagKinds(p.Tags)

	return
}

func (s projectService) Fork(cmd *ProjectForkCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	name, err := s.genForkedRepoName(cmd.Owner, cmd.From.Name)
	if err != nil {
		return
	}

	pid, err := pr.Fork(cmd.From.RepoId, name)
	if err != nil {
		return
	}

	v := new(domain.Project)
	cmd.toProject(name, v)
	v.RepoId = pid

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// create activity
	ua := genActivityForCreatingResource(
		p.Owner, domain.ResourceTypeProject, p.Id,
	)

	_ = s.activity.Save(&ua)

	// send event
	_ = s.sender.IncreaseFork(&domain.ResourceIndex{
		Owner: cmd.From.Owner,
		Id:    cmd.From.Id,
	})

	return
}

func (s projectService) genForkedRepoName(
	owner domain.Account, srcName domain.ProjName,
) (domain.ProjName, error) {
	v, err := s.repo.ListAndSortByUpdateTime(
		owner,
		&repository.ResourceListOption{Name: srcName.ProjName()},
	)
	if err != nil {
		return nil, err
	}

	items := v.Projects
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
