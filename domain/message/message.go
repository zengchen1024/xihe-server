package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Sender interface {
	AddFollowing(msg domain.Following) error
	RemoveFollowing(msg domain.Following) error
}

type FollowingHandler interface {
	HandleEventAddFollowing(domain.Following) error
	HandleEventRemoveFollowing(msg domain.Following) error
}
