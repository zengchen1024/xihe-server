package controller

type RepoFileCreateRequest struct {
	Content string `json:"content"`
}

type RepoFileUpdateRequest = RepoFileCreateRequest
