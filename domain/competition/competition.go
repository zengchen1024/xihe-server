package competition

import "io"

type Competition interface {
	UploadSubmissionFile(data io.Reader, path string) error
}
