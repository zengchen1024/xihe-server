package domain

type UserLike struct {
	Owner Account

	Like
}

type Like struct {
	ResourceOwner Account
	ResourceType  ResourceType
	ResourceId    string
}
