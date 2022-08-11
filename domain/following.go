package domain

type Following struct {
	Owner    Account
	Account  Account
	AvatarId AvatarId
	Bio      Bio
}

type Follower struct {
	Owner    Account
	Account  Account
	AvatarId AvatarId
	Bio      Bio
}
