package domain

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
