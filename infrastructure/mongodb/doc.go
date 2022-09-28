package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId             = "id"
	fieldBio            = "bio"
	fieldDesc           = "desc"
	fieldCoverId        = "cover_id"
	fieldTags           = "tags"
	fieldName           = "name"
	fieldItems          = "items"
	fieldOwner          = "owner"
	fieldEmail          = "email"
	fieldAccount        = "account"
	fieldVersion        = "version"
	fieldRepoType       = "repo_type"
	fieldAvatarId       = "avatar_id"
	fieldFollower       = "follower"
	fieldFollowing      = "following"
	fieldLikeCount      = "like_count"
	fieldForkCount      = "fork_count"
	fieldIsFollower     = "is_follower"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldCreatedAt      = "created_at"
	fieldType           = "type"
	fieldModels         = "models"
	fieldDatasets       = "datasets"
	fieldRId            = "rid"
	fieldROwner         = "rowner"
	fieldRType          = "rtype"
	fieldUpdatedAt      = "updated_at"
	fieldDownloadCount  = "download_count"
	fieldFirstLetter    = "fl"
)

type dProject struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []projectItem `bson:"items" json:"-"`
}

type projectItem struct {
	Id        string `bson:"id"         json:"id"`
	Type      string `bson:"type"       json:"type"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	Training  string `bson:"training"   json:"training"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	ProjectPropertyItem `bson:",inline"`

	// These two items are not allowd to be set,
	// So, don't marshal it to avoid setting it occasionally.
	RelatedModels   []ResourceIndex `bson:"models"   json:"-"`
	RelatedDatasets []ResourceIndex `bson:"datasets" json:"-"`

	// Version, LikeCount and ForkCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"           json:"-"`
	LikeCount     int `bson:"like_count"        json:"-"`
	ForkCount     int `bson:"fork_count"        json:"-"`
	DownloadCount int `bson:"download_count"    json:"-"`
}

type ProjectPropertyItem struct {
	Name     string   `bson:"name"       json:"name"`
	FL       byte     `bson:"fl"         json:"fl"`
	Desc     string   `bson:"desc"       json:"desc"`
	CoverId  string   `bson:"cover_id"   json:"cover_id"`
	RepoType string   `bson:"repo_type"  json:"repo_type"`
	Tags     []string `bson:"tags"       json:"tags"`
}

type dModel struct {
	Owner string      `bson:"owner" json:"owner"`
	Items []modelItem `bson:"items" json:"-"`
}

type modelItem struct {
	Id        string   `bson:"id"         json:"id"`
	Name      string   `bson:"name"       json:"name"`
	Desc      string   `bson:"desc"       json:"desc"`
	Protocol  string   `bson:"protocol"   json:"protocol"`
	RepoType  string   `bson:"repo_type"  json:"repo_type"`
	RepoId    string   `bson:"repo_id"    json:"repo_id"`
	Tags      []string `bson:"tags"       json:"tags"`
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at"`

	// RelatedDatasets is not allowd to be set,
	// So, don't marshal it to avoid setting it occasionally.
	RelatedDatasets []ResourceIndex `bson:"datasets" json:"-"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version   int `bson:"version"       json:"-"`
	LikeCount int `bson:"like_count"    json:"-"`
}

type dDataset struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []datasetItem `bson:"items" json:"-"`
}

type datasetItem struct {
	Id        string   `bson:"id"         json:"id"`
	Name      string   `bson:"name"       json:"name"`
	Desc      string   `bson:"desc"       json:"desc"`
	Protocol  string   `bson:"protocol"   json:"protocol"`
	RepoType  string   `bson:"repo_type"  json:"repo_type"`
	RepoId    string   `bson:"repo_id"    json:"repo_id"`
	Tags      []string `bson:"tags"       json:"tags"`
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version   int `bson:"version"       json:"-"`
	LikeCount int `bson:"like_count"    json:"-"`
}

type DUser struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name                    string `bson:"name"       json:"name"`
	Email                   string `bson:"email"      json:"email"`
	Bio                     string `bson:"bio"        json:"bio"`
	AvatarId                string `bson:"avatar_id"  json:"avatar_id"`
	PlatformToken           string `bson:"token"      json:"token"`
	PlatformUserId          string `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string `bson:"nid"        json:"nid"`

	Follower  []string `bson:"follower"   json:"-"`
	Following []string `bson:"following"  json:"-"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type dLogin struct {
	Account string `bson:"account"   json:"account"`
	Info    string `bson:"info"      json:"info"`
}

type dLike struct {
	Owner string     `bson:"owner" json:"owner"`
	Items []likeItem `bson:"items" json:"-"`
}

type likeItem struct {
	CreatedAt int64 `bson:"created_at" json:"created_at"`

	ResourceObject `bson:",inline"`
}

type dActivity struct {
	Owner string         `bson:"owner" json:"owner"`
	Items []activityItem `bson:"items" json:"-"`
}

type activityItem struct {
	Type string `bson:"type" json:"type"`
	Time int64  `bson:"time" json:"time"`

	ResourceObject `bson:",inline"`
}

type ResourceObject struct {
	Id    string `bson:"rid"     json:"rid"`
	Type  string `bson:"rtype"   json:"rtype"`
	Owner string `bson:"rowner"  json:"rowner"`
}

type ResourceIndex struct {
	Id    string `bson:"rid"     json:"rid"`
	Owner string `bson:"rowner"  json:"rowner"`
}

type dResourceTags struct {
	ResourceType string         `bson:"rtype"   json:"rtype"`
	Items        []dDomainTags  `bson:"items"   json:"items"`
	Orders       map[string]int `bson:"orders"   json:"orders"`
}

type dDomainTags struct {
	Domain string   `bson:"domain"   json:"domain"`
	Kind   string   `bson:"kind"     json:"kind"`
	Order  int      `bson:"order"    json:"order"`
	Tags   []string `bson:"tags"     json:"tags"`
}
