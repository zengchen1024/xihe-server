package platform

// errorTooManyFilesToDelete
type errorTooManyFilesToDelete struct {
	error
}

func NewErrorTooManyFilesToDelete(err error) errorTooManyFilesToDelete {
	return errorTooManyFilesToDelete{err}
}

// helper
func IsErrorTooManyFilesToDelete(err error) bool {
	_, ok := err.(errorTooManyFilesToDelete)

	return ok
}
