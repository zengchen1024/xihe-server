package uploader

import "io"

type SubmissionFileUploader interface {
	Upload(data io.Reader, path string) error
}
