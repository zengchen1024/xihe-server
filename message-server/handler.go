package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

var _ message.EventHandler = (*handler)(nil)

const sleepTime = 100 * time.Millisecond

type handler struct {
	log *logrus.Entry

	maxRetry int
	user     app.UserService
	model    app.ModelService
	dataset  app.DatasetService
	project  app.ProjectService
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
		switch like.Type.ResourceType() {

		case domain.ResourceProject:
			err = h.project.AddLike(like.Owner, like.Id)

		case domain.ResourceDataset:
			err = h.dataset.AddLike(like.Owner, like.Id)

		case domain.ResourceModel:
			err = h.model.AddLike(like.Owner, like.Id)
		}

		if err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of adding like for owner:%s, rid:%s, err:%v",
					like.Owner.Account(), like.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) HandleEventRemoveLike(like domain.Like) (err error) {
	return h.do(func() (err error) {
		switch like.Type.ResourceType() {
		case domain.ResourceProject:
			err = h.project.AddLike(like.Owner, like.Id)

		case domain.ResourceDataset:
			err = h.dataset.AddLike(like.Owner, like.Id)

		case domain.ResourceModel:
			err = h.model.AddLike(like.Owner, like.Id)
		}

		if err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of removing like for owner:%s, rid:%s, err:%v",
					like.Owner.Account(), like.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) HandleEventFork(index domain.ResourceIndex) error {
	return h.do(func() (err error) {
		if err = h.project.IncreaseFork(index); err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of fork for owner:%s, rid:%s, err:%v",
					index.Owner.Account(), index.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) do(f func() error) (err error) {
	if err = f(); err == nil {
		return
	}

	for i := 1; i < h.maxRetry; i++ {
		time.Sleep(sleepTime)

		if err = f(); err == nil {
			return
		}
	}

	return h.errMaxRetry(err)
}

func (h *handler) errMaxRetry(err error) error {
	return fmt.Errorf("exceed max retry num, last err:%v", err)
}

func isResourceNotExists(err error) bool {
	_, ok := err.(repository.ErrorResourceNotExists)

	return ok
}
