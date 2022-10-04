package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Sender interface {
	AddFollowing(msg domain.FollowerInfo) error
	RemoveFollowing(msg domain.FollowerInfo) error

	AddLike(*domain.ResourceObject) error
	RemoveLike(*domain.ResourceObject) error

	IncreaseFork(*domain.ResourceIndex) error
}

type EventHandler interface {
	FollowingHandler
	LikeHandler
	ForkHandler
}

type FollowingHandler interface {
	HandleEventAddFollowing(domain.FollowerInfo) error
	HandleEventRemoveFollowing(domain.FollowerInfo) error
}

type LikeHandler interface {
	HandleEventAddLike(*domain.ResourceObject) error
	HandleEventRemoveLike(*domain.ResourceObject) error
}

type ForkHandler interface {
	HandleEventFork(*domain.ResourceIndex) error
}
