package challengeimpl

type Config struct {
	AIQuestion              AIQuestion `json:"ai_question"                  required:"true"`
	Competitions            []string   `json:"competitions"                 required:"true"`
	CompetitionSuccessScore int        `json:"Competition_success_score"    required:"true"`
}

type AIQuestion struct {
	AIQuestionId             string `json:"ai_question_id"               required:"true"`
	QuestionPoolId           string `json:"question_pool_id"             required:"true"`
	Timeout                  int    `json:"timeout"                      required:"true"`
	RetryTimes               int    `json:"retry_times"                  required:"true"`
	ChoiceQuestionsNum       int    `json:"choice_questions_num"         required:"true"`
	ChoiceQuestionsCount     int    `json:"choice_questions_count"       required:"true"`
	ChoiceQuestionsScore     int    `json:"choice_questions_score"       required:"true"`
	CompletionQuestionsNum   int    `json:"completion_questions_num"     required:"true"`
	CompletionQuestionsCount int    `json:"completion_questions_count"   required:"true"`
	CompletionQuestionsScore int    `json:"completion_questions_score"   required:"true"`
}
