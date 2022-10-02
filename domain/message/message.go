package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Sender interface {
	AddFollowing(msg domain.FollowerInfo) error
	RemoveFollowing(msg domain.FollowerInfo) error

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
	HandleEventAddFollowing(domain.FollowerInfo) error
	HandleEventRemoveFollowing(domain.FollowerInfo) error
}

type LikeHandler interface {
	HandleEventAddLike(domain.Like) error
	HandleEventRemoveLike(msg domain.Like) error
}

type ForkHandler interface {
	HandleEventFork(domain.ResourceIndex) error
}
