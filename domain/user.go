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
}

type PlatformUser struct {
	Id          string
	NamespaceId string
}
