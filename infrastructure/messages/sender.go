package messages

import (
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

var topics Topics

type Topics struct {
	Like      string `json:"like"      required:"true"`
	Fork      string `json:"fork"      required:"true"`
	Following string `json:"following" required:"true"`
}

func NewMessageSender() message.Sender {
	return sender{}
}

type sender struct{}

// Following
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

	return s.send(topics.Following, &v)
}

// Like
func (s sender) AddLike(msg domain.Like) error {
	return s.sendLike(msg, actionAdd)
}

func (s sender) RemoveLike(msg domain.Like) error {
	return s.sendLike(msg, actionRemove)
}

func (s sender) sendLike(msg domain.Like, action string) error {
	v := msgLike{
		Action: action,
		Owner:  msg.Owner.Account(),
		Type:   msg.Type.ResourceType(),
		Id:     msg.Id,
	}

	return s.send(topics.Like, &v)
}

// Fork
func (s sender) IncreaseFork(msg domain.ResourceIndex) error {
	v := msgFork{
		Owner: msg.Owner.Account(),
		Id:    msg.Id,
	}

	return s.send(topics.Fork, &v)
}

func (s sender) send(topic string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kafka.Publish(topic, &mq.Message{
		Body: body,
	})
}
