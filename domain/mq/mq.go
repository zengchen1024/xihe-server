package mq

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Sender interface {
	AddFollowing(msg domain.Following) error
	RemoveFollowing(msg domain.Following) error
}

type FollowingHandler interface {
	AddFollowing(domain.Following) error
	RemoveFollowing(msg domain.Following) error
}
