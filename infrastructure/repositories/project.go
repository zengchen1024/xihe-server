package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewProjectRepository(mapper ProjectMapper) repository.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p *domain.Project) (r domain.Project, err error) {
	if p.Id != "" {
		return
	}

	do := ProjectBasicDO{
		Owner:     p.Owner,
		Name:      p.Name.ProjName(),
		Desc:      p.Desc.ProjDesc(),
		Type:      p.Type.RepoType(),
		CoverId:   p.CoverId.CoverId(),
		Protocol:  p.Protocol.ProtocolName(),
		Training:  p.Training.TrainingSDK(),
		Inference: p.Inference.InferenceSDK(),
	}

	v, err := impl.mapper.Create(do)
	if err == nil {
		r.Id = v
	}
	return
}

func (impl project) Get(pid string) (r domain.Project, err error) {
	return
}

type ProjectBasicDO struct {
	Owner     string
	Name      string
	Desc      string
	Type      string
	CoverId   string
	Protocol  string
	Training  string
	Inference string
	Tags      []string
}

type ProjectDO struct {
	ProjectBasicDO

	Tags []string

	LikeAccount int

	AccumulatedDownloads int
	RecentDownloads      map[string]int
}

type ProjectMapper interface {
	Create(ProjectBasicDO) (string, error)
	Update(string, ProjectBasicDO) (string, error)
	Get(string) (ProjectDO, error)
}
