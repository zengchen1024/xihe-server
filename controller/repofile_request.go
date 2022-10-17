package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type RepoFileCreateRequest struct {
	Content string `json:"content"`
}

func (req *RepoFileCreateRequest) toCmd(repoId, path string) (cmd app.RepoFileCreateCmd, err error) {
	if cmd.Path, err = domain.NewFilePath(path); err != nil {
		return
	}

	cmd.Content = &req.Content
	cmd.RepoId = repoId

	return
}

type RepoFileUpdateRequest = RepoFileCreateRequest
