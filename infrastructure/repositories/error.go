package repositories

import "github.com/opensourceways/xihe-server/domain/repository"

type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{err}
}

func convertError(err error) (out error) {
	switch err.(type) {
	case ErrorDuplicateCreating:
		out = repository.NewErrorDuplicateCreating(err)

	default:
		out = err
	}

	return
}
