package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Sender interface {
	AddFollowing(msg domain.Following) error
	RemoveFollowing(msg domain.Following) error

	AddLike(msg domain.Like) error
	RemoveLike(msg domain.Like) error

	IncreaseFork(msg domain.ResourceIndex) error
}

type EventHandler interface {
	FollowingHandler
	LikeHandler
	ForkHandler
}

type FollowingHandler interface {
	HandleEventAddFollowing(domain.Following) error
	HandleEventRemoveFollowing(msg domain.Following) error
}

type LikeHandler interface {
	HandleEventAddLike(domain.Like) error
	HandleEventRemoveLike(msg domain.Like) error
}

type ForkHandler interface {
	HandleEventFork(domain.ResourceIndex) error
}
