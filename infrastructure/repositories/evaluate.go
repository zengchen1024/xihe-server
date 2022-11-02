package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type EvaluateMapper interface {
	Insert(*EvaluateDO, int) (string, error)
	Get(*EvaluateIndexDO) (EvaluateSummaryDO, error)
	UpdateDetail(*EvaluateIndexDO, *EvaluateDetailDO) error
	List(*ResourceIndexDO, string) ([]EvaluateSummaryDO, int, error)
	GetStandardEvaluateParms(*EvaluateIndexDO) (StandardEvaluateParmsDO, error)
}

func NewEvaluateRepository(mapper EvaluateMapper) repository.Evaluate {
	return evaluate{mapper}
}

type evaluate struct {
	mapper EvaluateMapper
}

func (impl evaluate) Save(ut *domain.Evaluate, version int) (string, error) {
	if ut.Id != "" {
		return "", errors.New("must be a new project")
	}

	do := impl.toEvaluateDO(ut)

	v, err := impl.mapper.Insert(&do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl evaluate) FindInstance(info *domain.EvaluateIndex) (
	r repository.EvaluateSummary, err error,
) {
	index := impl.toEvaluateIndexDO(info)
	v, err := impl.mapper.Get(&index)
	if err != nil {
		err = convertError(err)

		return
	}

	r.EvaluateDetail = v.EvaluateDetailDO
	r.Id = v.Id

	return
}

func (impl evaluate) FindInstances(info *domain.ResourceIndex, lastCommit string) (
	r []repository.EvaluateSummary, version int, err error,
) {
	index := toResourceIndexDO(info)
	v, version, err := impl.mapper.List(&index, lastCommit)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]repository.EvaluateSummary, len(v))

	for i := range v {
		r[i].EvaluateDetail = v[i].EvaluateDetailDO
		r[i].Id = v[i].Id
	}

	return
}

func (impl evaluate) GetStandardEvaluateParms(info *domain.EvaluateIndex) (
	domain.StandardEvaluateParms, error,
) {
	index := impl.toEvaluateIndexDO(info)

	v, err := impl.mapper.GetStandardEvaluateParms(&index)
	if err != nil {
		err = convertError(err)

		return domain.StandardEvaluateParms{}, err
	}

	return v, nil
}

func (impl evaluate) UpdateDetail(
	info *domain.EvaluateIndex, detail *domain.EvaluateDetail,
) error {
	index := impl.toEvaluateIndexDO(info)

	if err := impl.mapper.UpdateDetail(&index, detail); err != nil {
		return convertError(err)
	}

	return nil
}
