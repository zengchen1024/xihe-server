package messages

import (
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

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
	v := msgLike{Action: action}

	toMsgResourceObject(msg, &v.Resource)

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

// Training
func (s sender) CreateTraining(info *domain.TrainingIndex) error {
	v := msgTraining{
		User:       info.Project.Owner.Account(),
		ProjectId:  info.Project.Id,
		TrainingId: info.TrainingId,
	}

	return s.send(topics.Training, &v)
}

// Finetune
func (s sender) CreateFinetune(info *domain.FinetuneIndex) error {
	v := msgFinetune{
		User: info.Owner.Account(),
		Id:   info.Id,
	}

	return s.send(topics.Finetune, &v)
}

// Inference
func (s sender) CreateInference(info *domain.InferenceInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionCreate
	v.ProjectName = info.ProjectName.ResourceName()

	return s.send(topics.Inference, &v)

}

func (s sender) ExtendInferenceSurvivalTime(info *message.InferenceExtendInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionExtend
	v.Expiry = info.Expiry

	return s.send(topics.Inference, &v)
}

func (s sender) toInferenceMsg(index *domain.InferenceIndex) msgInference {
	return msgInference{
		ProjectId:    index.Project.Id,
		LastCommit:   index.LastCommit,
		InferenceId:  index.Id,
		ProjectOwner: index.Project.Owner.Account(),
	}
}

// Evaluate
func (s sender) CreateEvaluate(info *message.EvaluateInfo) error {
	v := msgEvaluate{
		Type:         info.Type,
		OBSPath:      info.OBSPath,
		ProjectId:    info.Project.Id,
		TrainingId:   info.TrainingId,
		EvaluateId:   info.Id,
		ProjectOwner: info.Project.Owner.Account(),
	}

	return s.send(topics.Evaluate, &v)
}

// RelatedResource
func (s sender) AddRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionAdd)
}

func (s sender) RemoveRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionRemove)
}

func (s sender) RemoveRelatedResources(msg *message.RelatedResources) error {
	v := msgRelatedResources{Action: actionRemove}

	toMsgResourceObject(&msg.Promoter, &v.Promoter)

	v.Resources = make([]resourceObjects, len(msg.Resources))
	for i := range msg.Resources {
		toMsgResourceObjects(&msg.Resources[i], &v.Resources[i])
	}

	return s.send(topics.RelatedResource, &v)
}

func (s sender) sendRelatedResource(msg *message.RelatedResource, action string) error {
	v := msgRelatedResources{Action: action}

	toMsgResourceObject(msg.Promoter, &v.Promoter)

	v.Resources = []resourceObjects{
		{
			Type: msg.Resource.Type.ResourceType(),
			Objects: []resourceIndex{
				{
					Owner: msg.Resource.Owner.Account(),
					Id:    msg.Resource.Id,
				},
			},
		},
	}

	return s.send(topics.RelatedResource, &v)
}

// Competition
func (s sender) CalcScore(info *message.SubmissionInfo) error {
	v := msgSubmission{
		CId:   info.Index.Id,
		Phase: info.Index.Phase.CompetitionPhase(),
		SId:   info.Id,
		Path:  info.OBSPath,
	}

	return s.send(topics.Submission, &v)
}

// send
func (s sender) send(topic string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kafka.Publish(topic, &mq.Message{
		Body: body,
	})
}
