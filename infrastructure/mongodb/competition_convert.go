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
		Duration: doc.Duration,
	}
}

func (col competition) toCompetitionDO(
	doc *DCompetition, do *repositories.CompetitionDO,
) {
	do.Enabled = doc.Enabled
	do.Phase = doc.Phase
	do.Doc = doc.Doc
	do.Forum = doc.Forum
	do.DatasetDoc = doc.DatasetDoc
	do.DatasetURL = doc.DatasetURL

	col.toCompetitionSummaryDO(doc, &do.CompetitionSummaryDO)
}

func (col competition) toCompetitorDO(
	doc *dCompetitor, teamName string, do *repositories.CompetitorDO,
) {
	*do = repositories.CompetitorDO{
		CompetitorInfoDO: repositories.CompetitorInfoDO{
			Name:     doc.Name,
			City:     doc.City,
			Email:    doc.Email,
			Phone:    doc.Phone,
			Account:  doc.Account,
			Identity: doc.Identity,
			Province: doc.Province,
			Detail:   doc.Detail,
		},
		TeamId:   doc.TeamId,
		TeamName: teamName,
		TeamRole: doc.TeamRole,
	}
}

func toCompetitorInfoDOC(
	doc *DCompetitorInfo, do *repositories.CompetitorInfoDO,
) {
	*doc = DCompetitorInfo{
		Name:     do.Name,
		City:     do.City,
		Email:    do.Email,
		Phone:    do.Phone,
		Account:  do.Account,
		Identity: do.Identity,
		Province: do.Province,
		Detail:   do.Detail,
	}
}

func (col competition) toCompetitionSubmissionDO(
	doc *dSubmission, do *repositories.CompetitionSubmissionDO,
) {
	*do = repositories.CompetitionSubmissionDO{
		Id:         doc.Id,
		Status:     doc.Status,
		OBSPath:    doc.OBSPath,
		SubmitAt:   doc.SubmitAt,
		Score:      float32(doc.Score),
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

func (col competition) toCompetitionRepoDO(
	doc *dCompetitionRepo, do *repositories.CompetitionRepoDO,
) {
	*do = repositories.CompetitionRepoDO{
		TeamId:     doc.TeamId,
		Individual: doc.Individual,
		Owner:      doc.Owner,
		Repo:       doc.Repo,
	}
}

func (col competition) toSubmissionDoc(
	do *repositories.CompetitionSubmissionDO,
	doc *dSubmission, date string,
) {
	*doc = dSubmission{
		Id:         do.Id,
		TeamId:     do.TeamId,
		Individual: do.Individual,
		Status:     do.Status,
		OBSPath:    do.OBSPath,
		SubmitAt:   do.SubmitAt,
		Score:      float64(do.Score),
		Date:       date,
	}
}
