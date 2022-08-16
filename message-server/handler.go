package main

import (
	"fmt"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type handler struct {
	msxRetry int
	user     app.UserService
}

func (h *handler) HandleEventAddFollowing(f domain.Following) error {
	return h.do(func() (err error) {
		if err = h.user.AddFollower(f.Account, f.Owner); err == nil {
			return
		}

		if _, ok := err.(repository.ErrorDuplicateCreating); ok {
			err = nil
		}

		return
	})
}

func (h *handler) HandleEventRemoveFollowing(f domain.Following) (err error) {
	return h.do(func() error {
		return h.user.RemoveFollower(f.Account, f.Owner)
	})
}

func (h *handler) do(f func() error) (err error) {
	for i := 0; i < h.msxRetry; i++ {
		if err = f(); err == nil {
			return
		}
	}

	return h.errMaxRetry(err)
}

func (h *handler) errMaxRetry(err error) error {
	return fmt.Errorf("exceed max retry num, last err:%v", err)
}
