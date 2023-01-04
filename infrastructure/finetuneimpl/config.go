package finetuneimpl

type Config struct {
	Endpoint           string   `json:"endpoint"              required:"true"`
	JobDoneStatus      []string `json:"job_done_status"       required:"true"`
	CanTerminateStatus []string `json:"can_terminate_status"  required:"true"`
}
