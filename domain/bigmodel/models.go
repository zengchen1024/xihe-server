package bigmodel

import (
	"io"

	"github.com/opensourceways/xihe-server/domain"
)

type BigModel interface {
	DescribePicture(io.Reader, string, int64) (string, error)
	GenPicture(domain.Account, string) (string, error)
	GenPictures(domain.Account, string) ([]string, error)
	Ask(domain.Question, domain.OBSFile) (string, error)
}
