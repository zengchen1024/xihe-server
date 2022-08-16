package message

import (
	"context"
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	libmq "github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

func Subscribe(ctx context.Context, handler interface{}, log *logrus.Entry) error {
	subscribers := make(map[string]libmq.Subscriber)

	defer func() {
		for k, s := range subscribers {
			if err := s.Unsubscribe(); err != nil {
				log.Errorf("failed to unsubscribe for topic:%s, err:%v", k, err)
			}
		}
	}()

	if h, ok := handler.(message.FollowingHandler); ok {
		s, err := registerFollowingHandler(h)
		if err != nil {
			return err
		}

		subscribers[topicFollowing] = s
	}

	if len(subscribers) == 0 {
		return nil
	}

	<-ctx.Done()

	return nil
}

func registerFollowingHandler(h message.FollowingHandler) (libmq.Subscriber, error) {
	return kafka.Subscribe(topicFollowing, func(e libmq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgFollowing{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		f := domain.Following{}
		if f.Owner, err = domain.NewAccount(body.Owner); err != nil {
			return
		}

		if f.Account, err = domain.NewAccount(body.Following); err != nil {
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
