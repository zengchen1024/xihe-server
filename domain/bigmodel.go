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

type UserWuKongPicture struct {
	User Account
	WuKongPicture
}

type WuKongPicture struct {
	Id        string
	OBSPath   string
	CreatedAt string

	WuKongPictureMeta
}

type WuKongPictureInfo struct {
	WuKongPictureMeta

	Link string `json:"link"`
}

type WuKongPictureMeta struct {
	Style string            `json:"style"`
	Desc  WuKongPictureDesc `json:"desc"`
}
