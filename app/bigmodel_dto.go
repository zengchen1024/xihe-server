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

type WuKongCmd bigmodel.WuKongReq

func (cmd *WuKongCmd) Validate() error {
	if cmd.Desc == "" {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongPicturesListCmd = repository.WuKongPictureListOption

type WuKongPictureAddLikeCmd struct {
	User    domain.Account
	OBSPath string
}

type WuKongPictureDTO = bigmodel.WuKongPictureInfo

type WuKongPicturesDTO struct {
	Total    int                `json:"total"`
	Pictures []WuKongPictureDTO `json:"pictures"`
}

type UserLikedWuKongPictureDTO struct {
	WuKongPictureDTO

	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
}
