package domain

type UserLuoJiaRecord struct {
	User Account

	LuoJiaRecord
}

type LuoJiaRecord struct {
	Id        string
	CreatedAt int64
}

type LuoJiaRecordIndex struct {
	User Account
	Id   string
}
