package app

type ErrorExceedMaxRelatedResourceNum struct {
	error
}

type ErrorPrivateRepo struct {
	error
}

type ErrorExccedMaxTrainingRecord struct {
	error
}

type ErrorOnlyOneRunningTraining struct {
	error
}

type ErrorUnavailableRepoFile struct {
	error
}

type ErrorUpdateLFSFile struct {
	error
}

type ErrorPreviewLFSFile struct {
	error
}

type errorUnavailableTraining struct {
	error
}

func IsErrorUnavailableTraining(err error) bool {
	_, ok := err.(errorUnavailableTraining)

	return ok
}
