package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionMapper interface {
	List(*CompetitionListOptionDO) ([]CompetitionSummaryDO, error)
	Get(index *CompetitionIndexDO, competitor string) (CompetitionDO, bool, error)
	GetTeam(index *CompetitionIndexDO, competitor string) ([]CompetitorDO, error)
	GetResult(*CompetitionIndexDO) (
		bool, []CompetitionTeamDO, []CompetitionSubmissionDO, error,
	)
	GetSubmisstions(cid, competitor string) (
		CompetitionRepoDO, []CompetitionSubmissionDO, error,
	)
}

func NewCompetitionRepository(mapper CompetitionMapper) repository.Competition {
	return competition{mapper}
}

type competition struct {
	mapper CompetitionMapper
}

func (impl competition) List(opt *repository.CompetitionListOption) (
	[]repository.CompetitionSummary, error,
) {
	do := new(CompetitionListOptionDO)
	impl.toCompetitionListOptionDO(opt, do)

	v, err := impl.mapper.List(do)
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
		if err = v[i].toCompetitor(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl competition) GetResult(index *domain.CompetitionIndex) (
	order domain.CompetitionScoreOrder,
	teams []domain.CompetitionTeam,
	results []domain.CompetitionSubmission, err error,
) {

	do := impl.toCompetitionIndexDO(index)

	b, ts, rs, err := impl.mapper.GetResult(&do)
	if err != nil || len(rs) == 0 {
		return
	}

	order = domain.NewCompetitionScoreOrder(b)

	teams = make([]domain.CompetitionTeam, len(ts))
	for i := range ts {
		if err = ts[i].toTeam(&teams[i]); err != nil {
			return
		}
	}

	results = make([]domain.CompetitionSubmission, len(rs))
	for i := range rs {
		if err = rs[i].toCompetitionResult(&results[i]); err != nil {
			return
		}
	}

	return
}

func (impl competition) GetSubmisstions(cid string, c domain.Account) (
	repo domain.CompetitionRepo,
	results []domain.CompetitionSubmission, err error,
) {
	r, rs, err := impl.mapper.GetSubmisstions(cid, c.Account())
	if err != nil || len(rs) == 0 {
		return
	}

	results = make([]domain.CompetitionSubmission, len(rs))
	for i := range rs {
		if err = rs[i].toCompetitionResult(&results[i]); err != nil {
			return
		}
	}

	if r.Owner != "" {
		err = r.toCompetitionRepo(&repo)
	}

	return
}
