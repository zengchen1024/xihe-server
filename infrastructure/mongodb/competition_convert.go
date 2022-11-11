package mongodb

import "github.com/opensourceways/xihe-server/infrastructure/repositories"

func (col competition) toCompetitionSummaryDO(
	doc *DCompetition, do *repositories.CompetitionSummaryDO,
) {
	*do = repositories.CompetitionSummaryDO{
		Id:       doc.Id.Hex(),
		Name:     doc.Name,
		Desc:     doc.Desc,
		Host:     doc.Host,
		Bonus:    doc.Bonus,
		Status:   doc.Status,
		Poster:   doc.Poster,
		Duration: doc.DatasetURL,
	}
}

func (col competition) toCompetitionDO(
	doc *DCompetition, do *repositories.CompetitionDO,
) {
	do.Doc = doc.Doc
	do.DatasetDoc = doc.DatasetDoc
	do.DatasetURL = doc.DatasetURL

	col.toCompetitionSummaryDO(doc, &do.CompetitionSummaryDO)
}

func (col competition) toCompetitorDO(
	doc *dCompetitor, teamName string, do *repositories.CompetitorDO,
) {
	*do = repositories.CompetitorDO{
		Name:     doc.Name,
		City:     doc.City,
		Email:    doc.Email,
		Phone:    doc.Phone,
		Account:  doc.Account,
		Identity: doc.Identity,
		Province: doc.Province,
		Detail:   doc.Detail,
		TeamId:   doc.TeamId,
		TeamName: teamName,
		TeamRole: doc.TeamRole,
	}
}
