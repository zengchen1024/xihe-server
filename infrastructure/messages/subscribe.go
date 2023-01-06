package messages

import (
	"context"
	"encoding/json"
	"errors"

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

	// training
	s, err = registerHandlerForTraining(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// finetune
	s, err = registerHandlerForFinetune(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	} else {
		log.Infof("does not subscribe finetune")
	}

	// inference
	s, err = registerHandlerForInference(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// evaluate
	s, err = registerHandlerForEvaluate(handler)
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

		body := msgRelatedResources{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return body.handle(h.HandleEventAddRelatedResource)

		case actionRemove:
			return body.handle(h.HandleEventRemoveRelatedResource)
		}

		return nil
	})
}

func registerHandlerForTraining(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.TrainingHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Training, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgTraining{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		if body.ProjectId == "" || body.TrainingId == "" {
			err = errors.New("invalid message of training")

			return
		}

		v := domain.TrainingIndex{}

		if v.Project.Owner, err = domain.NewAccount(body.User); err != nil {
			return
		}

		v.Project.Id = body.ProjectId
		v.TrainingId = body.TrainingId

		return h.HandleEventCreateTraining(&v)
	})
}

func registerHandlerForFinetune(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.FinetuneHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Finetune, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgFinetune{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		if body.Id == "" {
			err = errors.New("invalid message of finetune")

			return
		}

		v := domain.FinetuneIndex{Id: body.Id}
		if v.Owner, err = domain.NewAccount(body.User); err != nil {
			return
		}

		return h.HandleEventCreateFinetune(&v)
	})
}

func registerHandlerForInference(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.InferenceHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Inference, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgInference{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		v := domain.InferenceIndex{}

		if v.Project.Owner, err = domain.NewAccount(body.ProjectOwner); err != nil {
			return
		}

		v.Id = body.InferenceId
		v.Project.Id = body.ProjectId
		v.LastCommit = body.LastCommit

		switch body.Action {
		case actionCreate:
			info := domain.InferenceInfo{
				InferenceIndex: v,
			}

			info.ProjectName, err = domain.NewResourceName(body.ProjectName)
			if err != nil {
				return
			}

			return h.HandleEventCreateInference(&info)

		case actionExtend:
			return h.HandleEventExtendInferenceSurvivalTime(
				&message.InferenceExtendInfo{
					InferenceIndex: v,
					Expiry:         body.Expiry,
				},
			)
		}

		return nil
	})
}

func registerHandlerForEvaluate(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.EvaluateHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Evaluate, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgEvaluate{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		v := message.EvaluateInfo{}

		if v.Project.Owner, err = domain.NewAccount(body.ProjectOwner); err != nil {
			return
		}

		v.Id = body.EvaluateId
		v.Type = body.Type
		v.OBSPath = body.OBSPath
		v.Project.Id = body.ProjectId
		v.TrainingId = body.TrainingId

		return h.HandleEventCreateEvaluate(&v)
	})
}
