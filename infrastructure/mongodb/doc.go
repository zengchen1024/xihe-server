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
	Owner string      `bson:"owner" json:"owner"`
	Items []modelItem `bson:"items" json:"-"`
}

type modelItem struct {
	Id       string   `bson:"id"        json:"id"`
	Name     string   `bson:"name"      json:"name"`
	Desc     string   `bson:"desc"      json:"desc"`
	Protocol string   `bson:"protocol"  json:"protocol"`
	RepoType string   `bson:"repo_type" json:"repo_type"`
	Tags     []string `bson:"tags"      json:"tags"`
}

type dDataset struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []datasetItem `bson:"items" json:"-"`
}

type datasetItem struct {
	Id       string   `bson:"id"        json:"id"`
	Name     string   `bson:"name"      json:"name"`
	Desc     string   `bson:"desc"      json:"desc"`
	Protocol string   `bson:"protocol"  json:"protocol"`
	RepoType string   `bson:"repo_type" json:"repo_type"`
	Tags     []string `bson:"tags"      json:"tags"`
}

type dUser struct {
	Bio         string `bson:"bio"           json:"bio"`
	Email       string `bson:"email"         json:"email"`
	Account     string `bson:"account"       json:"account"`
	Password    string `bson:"password"      json:"password"`
	Nickname    string `bson:"nickname"      json:"nickname"`
	AvatarId    string `bson:"avatar_id"     json:"avatar_id"`
	PhoneNumber string `bson:"phone_number"  json:"phone_number"`
}
