package domain

type UserBasicInfo struct {
	Nickname Nickname
	AvatarId AvatarId
	Bio      Bio
}

type User struct {
	Id          string
	Email       Email
	Account     Account
	PhoneNumber PhoneNumber

	UserBasicInfo
}
