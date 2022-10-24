package repository

type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{err}
}

type ErrorResourceNotExists struct {
	error
}

func NewErrorResourceNotExists(err error) ErrorResourceNotExists {
	return ErrorResourceNotExists{err}
}

type ErrorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{err}
}

func IsErrorResourceNotExists(err error) bool {
	_, ok := err.(ErrorResourceNotExists)

	return ok
}

func IsErrorDuplicateCreating(err error) bool {
	_, ok := err.(ErrorDuplicateCreating)

	return ok
}
