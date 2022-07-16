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
