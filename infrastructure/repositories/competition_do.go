package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
)

type CompetitionSummaryDO struct {
	Bonus    int
	Id       string
	Name     string
	Desc     string
	Host     string
	Status   string
	Poster   string
	Duration string

	CompetitorsCount int
}

func (do *CompetitionSummaryDO) toCompetitionSummary(
	c *domain.CompetitionSummary,
) (err error) {
	c.Id = do.Id

	if c.Bonus, err = domain.NewCompetitionBonus(do.Bonus); err != nil {
		return
	}

	if c.Name, err = domain.NewCompetitionName(do.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewCompetitionDesc(do.Desc); err != nil {
		return
	}

	if c.Host, err = domain.NewCompetitionHost(do.Host); err != nil {
		return
	}

	if c.Status, err = domain.NewCompetitionStatus(do.Status); err != nil {
		return
	}

	if c.Poster, err = domain.NewURL(do.Poster); err != nil {
		return
	}

	c.Duration, err = domain.NewCompetitionDuration(do.Duration)

	return
}

type CompetitionDO struct {
	CompetitionSummaryDO

	Doc        string
	DatasetDoc string
	DatasetURL string
}

func (do *CompetitionDO) toCompetition(
	c *domain.Competition,
) (err error) {
	if c.Doc, err = domain.NewURL(do.Poster); err != nil {
		return
	}

	if c.DatasetDoc, err = domain.NewURL(do.Poster); err != nil {
		return
	}

	if c.DatasetURL, err = domain.NewURL(do.Poster); err != nil {
		return
	}

	err = do.CompetitionSummaryDO.toCompetitionSummary(&c.CompetitionSummary)

	return
}
