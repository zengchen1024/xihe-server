package main

import (
	"fmt"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

var _ message.EventHandler = (*handler)(nil)

type handler struct {
	maxRetry int
	user     app.UserService
	project  app.ProjectService
	model    app.ModelService
	dataset  app.DatasetService
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

func (h *handler) HandleEventAddLike(like domain.Like) error {
	return h.do(func() (err error) {
		switch like.ResourceType.ResourceType() {

		case domain.ResourceProject:
			err = h.project.AddLike(like.ResourceOwner, like.ResourceId)

		case domain.ResourceDataset:
			err = h.dataset.AddLike(like.ResourceOwner, like.ResourceId)

		case domain.ResourceModel:
			err = h.model.AddLike(like.ResourceOwner, like.ResourceId)
		}

		if err != nil {
			// TODO err = nil if no account or no resource
		}

		return
	})
}

func (h *handler) HandleEventRemoveLike(like domain.Like) (err error) {
	return h.do(func() (err error) {
		switch like.ResourceType.ResourceType() {

		case domain.ResourceProject:
			err = h.project.AddLike(like.ResourceOwner, like.ResourceId)

		case domain.ResourceDataset:
			err = h.dataset.AddLike(like.ResourceOwner, like.ResourceId)

		case domain.ResourceModel:
			err = h.model.AddLike(like.ResourceOwner, like.ResourceId)
		}

		if err != nil {
			// TODO err = nil if no account or no resource
		}

		return
	})
}

func (h *handler) do(f func() error) (err error) {
	for i := 0; i < h.maxRetry; i++ {
		if err = f(); err == nil {
			return
		}
	}

	return h.errMaxRetry(err)
}

func (h *handler) errMaxRetry(err error) error {
	return fmt.Errorf("exceed max retry num, last err:%v", err)
}
