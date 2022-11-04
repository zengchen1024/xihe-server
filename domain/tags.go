package domain

type DomainTags struct {
	Name   string
	Domain string `json:"domain"`
	Items  []Tags `json:"items"`
}

type Tags struct {
	Kind  string   `json:"kind"`
	Items []string `json:"items"`
}
