package domain

type UserLike struct {
	Owner Account

	Like
}

type Like struct {
	/*
		CreatedAt int64

		ResourceObj
	*/
	ResourceOwner Account
	ResourceType  ResourceType
	ResourceId    string
}
