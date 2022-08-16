package main

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type handler struct {
	user app.UserService
}

func (h *handler) AddFollowing(f domain.Following) error {
	// TODO: retry if necessary
	return h.user.AddFollower(f.Account, f.Owner)
}

func (h *handler) RemoveFollowing(f domain.Following) error {
	// TODO: retry if necessary
	return h.user.RemoveFollower(f.Account, f.Owner)
}
