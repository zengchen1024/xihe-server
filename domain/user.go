package domain

type User struct {
	Id      string
	Email   Email
	Account Account

	Bio      Bio
	AvatarId AvatarId

	PlatformUser  PlatformUser
	PlatformToken string

	Version int

	// following fileds is not under the controlling of version
	FollowerCount  int
	FollowingCount int
}

type PlatformUser struct {
	Id          string
	NamespaceId string
}

type FollowerInfo struct {
	User     Account
	Follower Account
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
