package mq

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/domain"

	"github.com/opensourceways/community-robot-lib/kafka"
	libmq "github.com/opensourceways/community-robot-lib/mq"
)

const (
	topicFollowing = "following"
)

type sender struct{}

func (s sender) AddFollowing(msg domain.Following) error {
	v := msgFollowing{
		Action:    actionAdd,
		Owner:     msg.Owner.Account(),
		Following: msg.Account.Account(),
	}

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kafka.Publish(topicFollowing, &libmq.Message{
		Body: body,
	})
}

func (s sender) RemoveFollowing(msg domain.Following) error {
	v := msgFollowing{
		Action:    actionRemove,
		Owner:     msg.Owner.Account(),
		Following: msg.Account.Account(),
	}

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kafka.Publish(topicFollowing, &libmq.Message{
		Body: body,
	})
}
