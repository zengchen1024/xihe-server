package bigmodel

import (
	"io"

	"github.com/opensourceways/xihe-server/domain"
)

type CodeGeexReq struct {
	Lang    string
	Content string
}

type CodeGeexResp struct {
	Result string `json:"result"`
	Finish string `json:"finish"`
}

type BigModel interface {
	DescribePicture(io.Reader, string, int64) (string, error)
	GenPicture(domain.Account, string) (string, error)
	GenPictures(domain.Account, string) ([]string, error)
	Ask(domain.Question, string) (string, error)
	VQAUploadPicture(f io.Reader, u domain.Account, fileName string) error
	LuoJiaUploadPicture(f io.Reader, u domain.Account) error
	PanGu(string) (string, error)
	LuoJia(string) (string, error)
	CodeGeex(*CodeGeexReq) (CodeGeexResp, error)
	GetWuKongSampleId() string
	GenWuKongSampleNums(int) []int
}
