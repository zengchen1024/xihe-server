package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateScopeDO = domain.EvaluateScope
type StandardEvaluateParmsDO = domain.StandardEvaluateParms

type EvaluateIndexDO struct {
	Id         string
	TrainingId string
	Project    ResourceIndexDO
}

func (impl evaluate) toEvaluateIndexDO(index *domain.EvaluateIndex) EvaluateIndexDO {
	return EvaluateIndexDO{
		Id:         index.Id,
		TrainingId: index.TrainingId,
		Project:    toResourceIndexDO(&index.Project),
	}
}

type EvaluateDetailDO = domain.EvaluateDetail

type EvaluateSummaryDO struct {
	Id string

	EvaluateDetailDO
}

type EvaluateDO struct {
	Id           string
	Type         string
	ProjectId    string
	TrainingId   string
	ProjectOwner string

	Params StandardEvaluateParmsDO

	EvaluateDetailDO
}

func (impl evaluate) toEvaluateDO(obj *domain.Evaluate) EvaluateDO {
	return EvaluateDO{
		Id:           obj.Id,
		Type:         obj.EvaluateType,
		ProjectId:    obj.Project.Id,
		TrainingId:   obj.TrainingId,
		ProjectOwner: obj.Project.Owner.Account(),
		Params:       obj.StandardEvaluateParms,
	}
}
