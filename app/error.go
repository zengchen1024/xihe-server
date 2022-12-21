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

	// It exceed the max times for a day.
	ErrorCodeAIQuestionExceedMaxTimes           = "aiquestion_exceed_max_times"
	ErrorCodeAIQuestionSubmissionExpiry         = "aiquestion_submission_expiry"
	ErrorCodeAIQuestionSubmissionUnmatchedTimes = "aiquestion_submission_unmatched_times"

	ErrorRepoFileTooManyFilesToDelete = "repofile_too_many_files_to_delete"

	ErrorCompetitionDuplicateSubmission = "competition_duplicate_submission"

	ErrorBigModelSensitiveInfo = "bigmodel_sensitive_info"

	ErrorTrainNoLog        = "train_no_log"
	ErrorTrainNoOutput     = "train_no_output"
	ErrorTrainExccedMaxNum = "train_excced_max_num" // excced max training num for a user

	ErrorWuKongExccedMaxNum = "wukong_excced_max_num" // excced max num user can add like to picture
)
