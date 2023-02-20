package uploader

import "io"

type Uploader interface {
	UploadSubmissionFile(data io.Reader, path string) error
}
