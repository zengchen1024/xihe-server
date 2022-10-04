package messages

import (
	"context"
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

func Subscribe(ctx context.Context, handler interface{}, log *logrus.Entry) error {
	subscribers := make(map[string]mq.Subscriber)

	defer func() {
		for k, s := range subscribers {
			if err := s.Unsubscribe(); err != nil {
				log.Errorf("failed to unsubscribe for topic:%s, err:%v", k, err)
			}
		}
	}()

	// register following
	s, err := registerHandlerForFollowing(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register like
	s, err = registerHandlerForLike(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register fork
	s, err = registerHandlerForFork(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register related resource
	s, err = registerHandlerForRelatedResource(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register end
	if len(subscribers) == 0 {
		return nil
	}

	<-ctx.Done()

	return nil
}

func registerHandlerForFollowing(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.FollowingHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Following, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgFollower{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		f := &domain.FollowerInfo{}
		if f.User, err = domain.NewAccount(body.User); err != nil {
			return
		}

		if f.Follower, err = domain.NewAccount(body.Follower); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddFollowing(f)

		case actionRemove:
			return h.HandleEventRemoveFollowing(f)
		}

		return nil
	})
}

func registerHandlerForLike(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.LikeHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Like, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgLike{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		like := &domain.ResourceObject{}
		if err = body.Resource.toResourceObject(like); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddLike(like)

		case actionRemove:
			return h.HandleEventRemoveLike(like)
		}

		return nil
	})
}

func registerHandlerForFork(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.ForkHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Fork, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgFork{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		index := domain.ResourceIndex{}
		if index.Owner, err = domain.NewAccount(body.Owner); err != nil {
			return
		}

		index.Id = body.Id

		return h.HandleEventFork(&index)
	})
}

func registerHandlerForRelatedResource(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.RelatedResourceHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.RelatedResource, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgRelatedResource{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		promoter, resource := &domain.ResourceObject{}, &domain.ResourceObject{}
		if err = body.toResources(promoter, resource); err != nil {
			return
		}

		v := &message.RelatedResource{
			Promoter: promoter,
			Resource: resource,
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddRelatedResource(v)

		case actionRemove:
			return h.HandleEventRemoveRelatedResource(v)
		}

		return nil
	})
}
