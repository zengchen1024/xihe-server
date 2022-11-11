package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionMapper interface {
	List(status, phase string) ([]CompetitionSummaryDO, error)
	Get(index *CompetitionIndexDO, user string) (CompetitionDO, bool, error)
	GetTeam(index *CompetitionIndexDO, user string) ([]CompetitorDO, error)
	GetResult(*CompetitionIndexDO) (bool, []CompetitionTeamDO, []CompetitionResultDO, error)
}

func NewCompetitionRepository(mapper CompetitionMapper) repository.Competition {
	return competition{mapper}
}

type competition struct {
	mapper CompetitionMapper
}

func (impl competition) List(status domain.CompetitionStatus, phase domain.CompetitionPhase) (
	[]repository.CompetitionSummary, error,
) {
	s := ""
	if status != nil {
		s = status.CompetitionStatus()
	}

	v, err := impl.mapper.List(s, phase.CompetitionPhase())
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

func (impl competition) Get(index *domain.CompetitionIndex, user domain.Account) (
	r repository.CompetitionInfo, b bool, err error,
) {
	s := ""
	if user != nil {
		s = user.Account()
	}

	do := impl.toCompetitionIndexDO(index)
	v, b, err := impl.mapper.Get(&do, s)
	if err != nil {
		return
	}

	if err = v.toCompetition(&r.Competition); err != nil {
		return
	}

	r.CompetitorCount = v.CompetitorsCount

	return
}

func (impl competition) GetTeam(index *domain.CompetitionIndex, user domain.Account) (
	[]domain.Competitor, error,
) {
	do := impl.toCompetitionIndexDO(index)

	v, err := impl.mapper.GetTeam(&do, user.Account())
	if err != nil {
		return nil, err
	}

	r := make([]domain.Competitor, len(v))
	for i := range v {
		v[i].toCompetitor(&r[i])
	}

	return r, nil
}

func (impl competition) GetResult(index *domain.CompetitionIndex) (
	order domain.CompetitionScoreOrder,
	teams []domain.CompetitionTeam,
	results []domain.CompetitionResult, err error,
) {

	do := impl.toCompetitionIndexDO(index)

	b, ts, rs, err := impl.mapper.GetResult(&do)
	if err != nil || len(rs) == 0 {
		return
	}

	order = domain.NewCompetitionScoreOrder(b)

	teams = make([]domain.CompetitionTeam, len(ts))
	for i := range ts {
		ts[i].toTeam(&teams[i])
	}

	results = make([]domain.CompetitionResult, len(rs))
	for i := range rs {
		rs[i].toCompetitionResult(&results[i])
	}

	return
}
