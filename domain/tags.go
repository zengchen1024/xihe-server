package domain

type ResourceTags struct {
	Type  ResourceType
	Items []DomainTags
}

type DomainTags struct {
	Domain string `json:"domain"`
	Items  []Tags `json:"items"`
}

type Tags struct {
	Kind  string   `json:"kind"`
	Items []string `json:"items"`
}
