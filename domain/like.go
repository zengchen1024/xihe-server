package domain

type UserLike struct {
	Owner Account

	Like
}

type Like struct {
	CreatedAt string

	ResourceObj
}
