package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectMapper interface {
	Insert(ProjectDO) (string, error)
	Update(ProjectDO) error
	Get(string, string) (ProjectDO, error)
	GetByName(string, string) (ProjectDO, error)
	List(string, ResourceListDO) ([]ProjectDO, error)
	ListUsersProjects(map[string][]string) ([]ProjectDO, error)

	AddLike(string, string) error
	RemoveLike(string, string) error
}

func NewProjectRepository(mapper ProjectMapper) repository.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p *domain.Project) (r domain.Project, err error) {
	if p.Id != "" {
		if err = impl.mapper.Update(impl.toProjectDO(p)); err != nil {
			err = convertError(err)
		} else {
			r = *p
			r.Version += 1
		}

		return
	}

	v, err := impl.mapper.Insert(impl.toProjectDO(p))
	if err != nil {
		err = convertError(err)
	} else {
		r = *p
		r.Id = v
	}

	return
}

func (impl project) Get(owner domain.Account, identity string) (r domain.Project, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) GetByName(owner domain.Account, name domain.ProjName) (
	r domain.Project, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.ProjName())
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) List(owner domain.Account, option repository.ResourceListOption) (
	r []domain.Project, err error,
) {
	do := ResourceListDO{
		Name: option.Name,
	}
	if option.RepoType != nil {
		do.RepoType = option.RepoType.RepoType()
	}

	v, err := impl.mapper.List(owner.Account(), do)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]domain.Project, len(v))
	for i := range v {
		if err = v[i].toProject(&r[i]); err != nil {
			return
		}
	}

	return
}

func (impl project) FindUserProjects(opts []repository.UserResourceListOption) (
	[]domain.Project, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersProjects(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.Project, len(v))
	for i := range v {
		if err = v[i].toProject(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl project) toProjectDO(p *domain.Project) ProjectDO {
	return ProjectDO{
		Id:       p.Id,
		Owner:    p.Owner.Account(),
		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ResourceDesc(),
		Type:     p.Type.ProjType(),
		CoverId:  p.CoverId.CoverId(),
		RepoType: p.RepoType.RepoType(),
		Protocol: p.Protocol.ProtocolName(),
		Training: p.Training.TrainingPlatform(),
		Tags:     p.Tags,
		RepoId:   p.RepoId,
		Version:  p.Version,
	}
}

type ProjectListDO struct {
	Name     string
	RepoType string
}

type ProjectDO struct {
	Id        string
	Owner     string
	Name      string
	Desc      string
	Type      string
	CoverId   string
	Protocol  string
	Training  string
	RepoType  string
	RepoId    string
	Tags      []string
	Version   int
	LikeCount int
}

func (do *ProjectDO) toProject(r *domain.Project) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewProjName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	if r.Type, err = domain.NewProjType(do.Type); err != nil {
		return
	}

	if r.CoverId, err = domain.NewConverId(do.CoverId); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	if r.Training, err = domain.NewTrainingPlatform(do.Training); err != nil {
		return
	}

	r.RepoId = do.RepoId
	r.Tags = do.Tags
	r.Version = do.Version
	r.LikeCount = do.LikeCount

	return
}
