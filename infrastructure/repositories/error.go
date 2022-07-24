package repositories

import "github.com/opensourceways/xihe-server/domain/repository"

type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{err}
}

type ErrorDataNotExists struct {
	error
}

func NewErrorDataNotExists(err error) ErrorDataNotExists {
	return ErrorDataNotExists{err}
}

type ErrorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{err}
}

func convertError(err error) (out error) {
	switch err.(type) {
	case ErrorDuplicateCreating:
		out = repository.NewErrorDuplicateCreating(err)

	case ErrorDataNotExists:
		out = repository.NewErrorResourceNotExists(err)

	case ErrorConcurrentUpdating:
		out = repository.NewErrorConcurrentUpdating(err)

	default:
		out = err
	}

	return
}
