package app

type CompetitionSummaryDTO struct {
	CompetitorCount int    `json:"count"`
	Bonus           int    `json:"bonus"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Host            string `json:"host"`
	Desc            string `json:"desc"`
	Status          string `json:"status"`
	Poster          string `json:"poster"`
	Duration        string `json:"duration"`
}

type CompetitionDTO struct {
	CompetitionSummaryDTO

	Doc                 string `json:"doc"`
	DatasetDoc          string `json:"dataset_doc"`
	DatasetDownloadAddr string `json:"dataset_download_addr"`
}

// ranking
type RankingDTO struct {
	Score    float32 `json:"score"`
	TeamName string  `json:"team_name"`
	SubmitAt string  `json:"submit_at"`
}

// team
type CompetitionTeamDTO struct {
	Name    string                     `json:"name"`
	Members []CompetitionTeamMemberDTO `json:"members"`
}

type CompetitionTeamMemberDTO struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

// result
type CompetitionResultDTO struct {
	RelatedProject string                       `json:"project"`
	Details        []CompetitionResultDetailDTO `json:"details"`
}

type CompetitionResultDetailDTO struct {
	SubmitAt string  `json:"submit_at"`
	FileName string  `json:"project"`
	Status   string  `json:"status"`
	Score    float32 `json:"score"`
}
