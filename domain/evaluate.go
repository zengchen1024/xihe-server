package domain

const (
	EvaluateTypeCustom   = "custom"
	EvaluateTypeStandard = "standard"
)

type EvaluateDetail = InferenceDetail

type Evaluate struct {
	EvaluateIndex

	StandardEvaluateParms

	EvaluateType string

	EvaluateDetail
}

type StandardEvaluateParms struct {
	MomentumScope     EvaluateScope
	BatchSizeScope    EvaluateScope
	LearningRateScope EvaluateScope
}

type EvaluateScope []string

type EvaluateIndex struct {
	TrainingIndex

	Id string
}
