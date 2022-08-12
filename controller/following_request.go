package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type followingCreateRequest struct {
	Account string `json:"account" required:"true"`
}

func (req *followingCreateRequest) toCmd() (cmd app.FollowingCreateCmd, err error) {
	cmd.Account, err = domain.NewAccount(req.Account)

	return
}
