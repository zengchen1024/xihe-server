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
	Like            string `json:"like"      required:"true"`
	Fork            string `json:"fork"      required:"true"`
	Following       string `json:"following" required:"true"`
	RelatedResource string `json:"related_resource" required:"true"`
}

func NewMessageSender() message.Sender {
	return sender{}
}

type sender struct{}

// Following
func (s sender) AddFollowing(msg *domain.FollowerInfo) error {
	return s.sendFollowing(msg, actionAdd)
}

func (s sender) RemoveFollowing(msg *domain.FollowerInfo) error {
	return s.sendFollowing(msg, actionRemove)
}

func (s sender) sendFollowing(msg *domain.FollowerInfo, action string) error {
	v := msgFollower{
		Action:   action,
		User:     msg.User.Account(),
		Follower: msg.Follower.Account(),
	}

	return s.send(topics.Following, &v)
}

// Like
func (s sender) AddLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionAdd)
}

func (s sender) RemoveLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionRemove)
}

func (s sender) sendLike(msg *domain.ResourceObject, action string) error {
	v := msgLike{
		Action:   action,
		Resource: toMsgResourceObject(msg),
	}

	return s.send(topics.Like, &v)
}

// Fork
func (s sender) IncreaseFork(msg *domain.ResourceIndex) error {
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

// RelatedResource
func (s sender) AddRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionAdd)
}

func (s sender) RemoveRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionRemove)
}

func (s sender) sendRelatedResource(msg *message.RelatedResource, action string) error {
	v := msgRelatedResource{
		Action:   action,
		Promoter: toMsgResourceObject(msg.Promoter),
		Resource: toMsgResourceObject(msg.Resource),
	}

	return s.send(topics.RelatedResource, &v)
}
