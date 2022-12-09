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

	// The current AI question is in-progress, but new request for
	// a new start comes.
	ErrorCodeChallengeInProgress = "challenge_in-progress"

	// It exceed the max times for a day.
	ErrorCodeChallengeExceedMaxTimes = "challenge_exceed_max_time"

	ErrorRepoFileTooManyFilesToDelete = "repofile_too_many_files_to_delete"

	ErrorCompetitionDuplicateSubmission = "competition_duplicate_submission"
)
