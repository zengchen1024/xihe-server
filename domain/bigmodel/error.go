package bigmodel

// errorSensitiveInfo
type errorSensitiveInfo struct {
	error
}

func NewErrorSensitiveInfo(err error) errorSensitiveInfo {
	return errorSensitiveInfo{err}
}

// helper
func IsErrorSensitiveInfo(err error) bool {
	_, ok := err.(errorSensitiveInfo)

	return ok
}
