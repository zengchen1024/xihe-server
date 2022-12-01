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

const (
	ErrorCodeSytem = "system"

	// The current ai question is in-progress, but new request for
	// a new start is comming.
	ErrorCodeChallengeInProgress = "challenge_in-progress"

	// It exceed the max times for a day.
	ErrorCodeChallengeExccedMaxTimes = "challenge_excced_max_time"
)
