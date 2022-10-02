package messages

const (
	actionAdd    = "add"
	actionRemove = "remove"
)

type msgFollower struct {
	Action   string `json:"action"`
	User     string `json:"user"`
	Follower string `json:"follower"`
}

type msgLike struct {
	Action string `json:"action"`
	Owner  string `json:"owner"`
	Type   string `json:"type"`
	Id     string `json:"id"`
}

type msgFork struct {
	Owner string `json:"owner"`
	Id    string `json:"id"`
}
