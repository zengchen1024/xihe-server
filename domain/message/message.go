package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateInfo struct {
	domain.EvaluateIndex
	Type    string
	OBSPath string
}

type Sender interface {
	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error

	AddLike(*domain.ResourceObject) error
	RemoveLike(*domain.ResourceObject) error

	IncreaseFork(*domain.ResourceIndex) error

	AddRelatedResource(*RelatedResource) error
	RemoveRelatedResource(*RelatedResource) error

	CreateTraining(*domain.TrainingInfo) error

	CreateInference(*domain.InferenceInfo) error
	ExtendInferenceExpiry(*domain.InferenceInfo) error

	CreateEvaluate(*EvaluateInfo) error
}

type EventHandler interface {
	RelatedResourceHandler
	FollowingHandler
	LikeHandler
	ForkHandler
}

type FollowingHandler interface {
	HandleEventAddFollowing(*domain.FollowerInfo) error
	HandleEventRemoveFollowing(*domain.FollowerInfo) error
}

type LikeHandler interface {
	HandleEventAddLike(*domain.ResourceObject) error
	HandleEventRemoveLike(*domain.ResourceObject) error
}

type ForkHandler interface {
	HandleEventFork(*domain.ResourceIndex) error
}

type RelatedResourceHandler interface {
	HandleEventAddRelatedResource(*RelatedResource) error
	HandleEventRemoveRelatedResource(*RelatedResource) error
}

type RelatedResource struct {
	Promoter *domain.ResourceObject
	Resource *domain.ResourceObject
}

type TrainingHandler interface {
	HandleEventCreateTraining(*domain.TrainingInfo) error
}

type InferenceHandler interface {
	HandleEventCreateInference(*domain.InferenceInfo) error
	HandleEventExtendInferenceExpiry(*domain.InferenceInfo) error
}

type EvaluateHandler interface {
	HandleEventCreateEvaluate(*EvaluateInfo) error
}
