package challengeimpl

type Config struct {
	AIQuestion               string   `json:"ai_question"                  required:"true"`
	Competitions             []string `json:"competitions"                 required:"true"`
	CompetitionSuccessStatus string   `json:"Competition_success_status"   required:"true"`
	CompetitionSuccessScore  int      `json:"Competition_success_score"    required:"true"`
	ChoiceQuestionsNum       int      `json:"choice_questions_num"         required:"true"`
	ChoiceQuestionsCount     int      `json:"choice_questions_count"       required:"true"`
	CompletionQuestionsNum   int      `json:"completion_questions_num"     required:"true"`
	CompletionQuestionsCount int      `json:"completion_questions_count"   required:"true"`
}
