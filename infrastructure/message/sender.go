package message

import (
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	libmq "github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

const (
	topicFollowing = "following"
)

func NewMessageSender() message.Sender {
	return sender{}
}

type sender struct{}

func (s sender) AddFollowing(msg domain.Following) error {
	return s.sendFollowing(msg, actionAdd)
}

func (s sender) RemoveFollowing(msg domain.Following) error {
	return s.sendFollowing(msg, actionRemove)
}

func (s sender) sendFollowing(msg domain.Following, action string) error {
	v := msgFollowing{
		Action:    action,
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
