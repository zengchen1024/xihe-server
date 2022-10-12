package main

import (
	"errors"
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

type likeHanler interface {
	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error
}

type relatedResourceHanler struct {
	Add    func(*domain.ReverselyRelatedResourceInfo) error
	Remove func(*domain.ReverselyRelatedResourceInfo) error
}

type handler struct {
	log *logrus.Entry

	maxRetry         int
	trainingEndpoint string

	user     app.UserService
	model    app.ModelService
	dataset  app.DatasetService
	project  app.ProjectService
	training app.TrainingService
}

func (h *handler) HandleEventAddRelatedResource(info *message.RelatedResource) error {
	rh := h.getHandlerForEventRelatedResource(info)
	if rh.Add == nil {
		return errors.New("unknown Reversely Related Resource")
	}

	data := h.getParameterForEventRelatedResource(info)

	return h.do(func() (err error) {
		return rh.Add(&data)
	})
}

func (h *handler) HandleEventRemoveRelatedResource(info *message.RelatedResource) error {
	rh := h.getHandlerForEventRelatedResource(info)
	if rh.Remove == nil {
		return errors.New("unknown Reversely Related Resource")
	}

	data := h.getParameterForEventRelatedResource(info)

	return h.do(func() (err error) {
		return rh.Remove(&data)
	})
}

func (h *handler) getParameterForEventRelatedResource(
	info *message.RelatedResource,
) domain.ReverselyRelatedResourceInfo {
	return domain.ReverselyRelatedResourceInfo{
		Promoter: &info.Promoter.ResourceIndex,
		Resource: &info.Resource.ResourceIndex,
	}
}

func (h *handler) getHandlerForEventRelatedResource(
	info *message.RelatedResource,
) (v relatedResourceHanler) {
	pt := info.Promoter.Type.ResourceType()

	switch info.Resource.Type.ResourceType() {
	case domain.ResourceTypeDataset.ResourceType():
		switch pt {
		case domain.ResourceTypeModel.ResourceType():
			v.Add = h.dataset.AddRelatedModel
			v.Remove = h.dataset.RemoveRelatedModel

		case domain.ResourceTypeProject.ResourceType():
			v.Add = h.dataset.AddRelatedProject
			v.Remove = h.dataset.RemoveRelatedProject
		}

	case domain.ResourceTypeModel.ResourceType():
		if pt == domain.ResourceTypeProject.ResourceType() {
			v.Add = h.model.AddRelatedProject
			v.Remove = h.model.RemoveRelatedProject
		}
	}

	return
}

func (h *handler) HandleEventAddFollowing(f *domain.FollowerInfo) error {
	return h.do(func() (err error) {
		if err = h.user.AddFollower(f); err == nil {
			return
		}

		if _, ok := err.(repository.ErrorDuplicateCreating); ok {
			err = nil
		}

		return
	})
}

func (h *handler) HandleEventRemoveFollowing(f *domain.FollowerInfo) (err error) {
	return h.do(func() error {
		return h.user.RemoveFollower(f)
	})
}

func (h *handler) HandleEventAddLike(obj *domain.ResourceObject) error {
	lh := h.getHandlerForEventLike(obj.Type)

	return h.handleEventLike(obj, "adding", lh.AddLike)
}

func (h *handler) HandleEventRemoveLike(obj *domain.ResourceObject) (err error) {
	lh := h.getHandlerForEventLike(obj.Type)

	return h.handleEventLike(obj, "removing", lh.RemoveLike)
}

func (h *handler) handleEventLike(
	obj *domain.ResourceObject, op string,
	f func(*domain.ResourceIndex) error,
) (err error) {
	return h.do(func() (err error) {
		if err = f(&obj.ResourceIndex); err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of %s like for owner:%s, rid:%s, err:%v",
					op, obj.Owner.Account(), obj.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) getHandlerForEventLike(t domain.ResourceType) likeHanler {
	switch t.ResourceType() {
	case domain.ResourceTypeProject.ResourceType():
		return h.project

	case domain.ResourceTypeDataset.ResourceType():
		return h.dataset

	case domain.ResourceTypeModel.ResourceType():
		return h.model
	}

	return nil
}

func (h *handler) HandleEventFork(index *domain.ResourceIndex) error {
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

func (h *handler) HandleEventCreateTraining(info *domain.TrainingInfo) error {
	return h.do(func() error {
		return h.training.CreateTrainingJob(info, h.trainingEndpoint)
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
