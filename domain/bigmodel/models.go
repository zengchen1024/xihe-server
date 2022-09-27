package bigmodel

import "io"

type BigModel interface {
	DescribePicture(io.Reader, string) (string, error)
}
