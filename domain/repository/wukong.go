package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type WuKongPictureListOption struct {
	CountPerPage int
	PageNum      int
}

type WuKongPictures struct {
	Pictures []string `json:"pictures"`
	Total    int      `json:"total"`
}

type WuKong interface {
	ListSamples(string, []int) ([]string, error)
	ListPictures(string, *WuKongPictureListOption) (WuKongPictures, error)
}

type WuKongPicture interface {
	List(user domain.Account) ([]domain.WuKongPicture, int, error)
	Save(*domain.UserWuKongPicture, int) (string, error)
	Delete(user domain.Account, pid string) error
	Get(user domain.Account, pid string) (domain.WuKongPicture, error)
}
