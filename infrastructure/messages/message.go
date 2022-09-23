package messages

const (
	actionAdd    = "add"
	actionRemove = "remove"
)

type msgFollowing struct {
	Action    string `json:"action"`
	Owner     string `json:"owner"`
	Following string `json:"following"`
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
