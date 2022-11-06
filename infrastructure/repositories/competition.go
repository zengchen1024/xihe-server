package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionMapper interface {
	List(status string) ([]CompetitionSummaryDO, error)
	Get(cid, user string) (CompetitionDO, bool, error)
}

func NewCompetitionRepository(mapper CompetitionMapper) repository.Competition {
	return competition{mapper}
}

type competition struct {
	mapper CompetitionMapper
}

func (impl competition) List(status domain.CompetitionStatus) (
	[]repository.CompetitionSummary, error,
) {
	s := ""
	if status != nil {
		s = status.CompetitionStatus()
	}

	v, err := impl.mapper.List(s)
	if err != nil {
		return nil, convertError(err)
	}

	if len(v) == 0 {
		return nil, err
	}

	r := make([]repository.CompetitionSummary, len(v))

	for i := range v {
		err := v[i].toCompetitionSummary(&r[i].CompetitionSummary)
		if err != nil {
			return nil, err
		}

		r[i].CompetitorCount = v[i].CompetitorsCount
	}

	return r, nil
}

func (impl competition) Get(cid string, user domain.Account) (
	r repository.CompetitionInfo, b bool, err error,
) {
	s := ""
	if user != nil {
		s = user.Account()
	}

	v, b, err := impl.mapper.Get(cid, s)
	if err != nil {
		return
	}

	if err = v.toCompetition(&r.Competition); err != nil {
		return
	}

	r.CompetitorCount = v.CompetitorsCount

	return
}
