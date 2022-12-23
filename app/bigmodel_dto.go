package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CodeGeexDTO = bigmodel.CodeGeexResp

type CodeGeexCmd bigmodel.CodeGeexReq

func (cmd *CodeGeexCmd) Validate() error {
	if cmd.Content == "" || cmd.Lang == "" {
		return errors.New("invalid cmd")
	}

	return nil
}

type LuoJiaRecordDTO struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
}

type WuKongCmd domain.WuKongPictureMeta

func (cmd *WuKongCmd) Validate() error {
	if cmd.Desc == nil {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongPicturesListCmd = repository.WuKongPictureListOption

type WuKongPictureAddLikeCmd struct {
	User    domain.Account
	OBSPath string
}

type WuKongPicturesDTO struct {
	Pictures []WuKongPictureInfoDTO `json:"pictures"`
	Total    int                    `json:"total"`
}

type UserLikedWuKongPictureDTO struct {
	WuKongPictureInfoDTO

	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

type WuKongPictureInfoDTO struct {
	Link  string `json:"link"`
	Desc  string `json:"desc"`
	Style string `json:"style"`
}
