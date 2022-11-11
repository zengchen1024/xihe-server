package mongodb

import "github.com/opensourceways/xihe-server/infrastructure/repositories"

func (col competition) toCompetitionSummaryDO(
	doc *DCompetition, do *repositories.CompetitionSummaryDO,
) {
	*do = repositories.CompetitionSummaryDO{
		Id:       doc.Id,
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

func (col competition) toCompetitionResultDO(
	doc *dSubmission, do *repositories.CompetitionResultDO,
) {
	*do = repositories.CompetitionResultDO{
		Id:         doc.Id,
		Status:     doc.Status,
		OBSPath:    doc.OBSPath,
		SubmitAt:   doc.SubmitAt,
		Score:      doc.Score,
		TeamId:     doc.TeamId,
		Individual: doc.Individual,
	}
}

func (col competition) toCompetitionTeamDO(
	doc *dTeam, do *repositories.CompetitionTeamDO,
) {
	*do = repositories.CompetitionTeamDO{
		Id:   doc.Id,
		Name: doc.Name,
	}
}
