package mongodb

import "github.com/opensourceways/xihe-server/infrastructure/repositories"

func (col competition) toCompetitionSummaryDO(
	doc *DCompetition, do *repositories.CompetitionSummaryDO,
) {
	do.Id = doc.Id
	do.Name = doc.Name
	do.Desc = doc.Desc
	do.Host = doc.Host
	do.Bonus = doc.Bonus
	do.Status = doc.Status
	do.Poster = doc.Poster
	do.Duration = doc.DatasetURL
}

func (col competition) toCompetitionDO(
	doc *DCompetition, do *repositories.CompetitionDO,
) {
	do.Doc = doc.Doc
	do.DatasetDoc = doc.DatasetDoc
	do.DatasetURL = doc.DatasetURL

	col.toCompetitionSummaryDO(doc, &do.CompetitionSummaryDO)
}
