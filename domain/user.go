package domain

type User struct {
	Id      string
	Email   Email
	Account Account

	Bio      Bio
	AvatarId AvatarId

	FollowerCount  int
	FollowingCount int

	PlatformUser  PlatformUser
	PlatformToken string

	Version int
}

type PlatformUser struct {
	Id          string
	NamespaceId string
}

type Following struct {
	Owner   Account
	Account Account
}

type Follower struct {
	Owner   Account
	Account Account
}

type FollowUserInfo struct {
	Account  Account
	AvatarId AvatarId
	Bio      Bio
}

type UserInfo struct {
	Account  Account
	AvatarId AvatarId
}
