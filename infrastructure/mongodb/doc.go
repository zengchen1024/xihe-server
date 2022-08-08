package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId       = "id"
	fieldName     = "name"
	fieldItems    = "items"
	fieldOwner    = "owner"
	fieldEmail    = "email"
	fieldVersion  = "version"
	fieldRepoType = "repo_type"
	fieldAccount  = "account"
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
	RepoId   string   `bson:"repo_id"   json:"repo_id"`
	Tags     []string `bson:"tags"      json:"tags"`
	Version  int      `bson:"version"    json:"-"`
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
	RepoId   string   `bson:"repo_id"   json:"repo_id"`
	Tags     []string `bson:"tags"      json:"tags"`
	Version  int      `bson:"version"    json:"-"`
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
	RepoId   string   `bson:"repo_id"   json:"repo_id"`
	Tags     []string `bson:"tags"      json:"tags"`
	Version  int      `bson:"version"    json:"-"`
}

type dUser struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name                    string `bson:"name"       json:"name"`
	Email                   string `bson:"email"      json:"email"`
	Bio                     string `bson:"bio"        json:"bio"`
	AvatarId                string `bson:"avatar_id"  json:"avatar_id"`
	PlatformToken           string `bson:"token"      json:"token"`
	PlatformUserId          string `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string `bson:"nid"        json:"nid"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type dLogin struct {
	Account string `bson:"account"   json:"account"`
	Info    string `bson:"info"      json:"info"`
}
