package mongodb

const (
	fieldId    = "id"
	fieldName  = "name"
	fieldItems = "items"
	fieldOwner = "owner"
)

type dProject struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []projectItem `bson:"items" json:"-"`
}

type projectItem struct {
	Id       string   `bson:"id"        json:"id"`
	Name     string   `bson:"name"      json:"name"`
	Desc     string   `bson:"desc"      json:"desc"`
	Type     string   `bson:"type"      json:"type"`
	CoverId  string   `bson:"cover_id"  json:"cover_id"`
	Protocol string   `bson:"protocol"  json:"protocol"`
	Training string   `bson:"training"  json:"training"`
	RepoType string   `bson:"repo_type" json:"repo_type"`
	Tags     []string `bson:"tags"      json:"tags"`
}

type dModel struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []projectItem `bson:"items" json:"-"`
}

type modelItem struct {
	Id       string   `bson:"id"        json:"id"`
	Name     string   `bson:"name"      json:"name"`
	Desc     string   `bson:"desc"      json:"desc"`
	Protocol string   `bson:"protocol"  json:"protocol"`
	RepoType string   `bson:"repo_type" json:"repo_type"`
	Tags     []string `bson:"tags"      json:"tags"`
}
