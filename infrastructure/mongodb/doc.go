package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId             = "id"
	fieldBio            = "bio"
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
	fieldIsFollower     = "is_follower"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldCreatedAt      = "created_at"
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

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version   int `bson:"version"       json:"-"`
	LikeCount int `bson:"like_count"    json:"-"`
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
	Id       string   `bson:"id"        json:"id"`
	Name     string   `bson:"name"      json:"name"`
	Desc     string   `bson:"desc"      json:"desc"`
	Protocol string   `bson:"protocol"  json:"protocol"`
	RepoType string   `bson:"repo_type" json:"repo_type"`
	RepoId   string   `bson:"repo_id"   json:"repo_id"`
	Tags     []string `bson:"tags"      json:"tags"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version   int `bson:"version"       json:"-"`
	LikeCount int `bson:"like_count"    json:"-"`
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

	Follower  []string `bson:"follower"   json:"-"`
	Following []string `bson:"following"  json:"-"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`

	/*
		// These 3 items are not the property of User
		// There are used to store the value when unmarshaling in method of GetByFollower
		IsFollower     bool `bson:"is_follower"       json:"-"`
		FollowerCount  int  `bson:"follower_count"    json:"-"`
		FollowingCount int  `bson:"following_count"   json:"-"`
	*/
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

	ResourceObj `bson:",inline"`
}

type dActivity struct {
	Owner string         `bson:"owner" json:"owner"`
	Items []activityItem `bson:"items" json:"-"`
}

type activityItem struct {
	Type string `bson:"type" json:"type"`
	Time int64  `bson:"time" json:"time"`

	ResourceObj `bson:",inline"`
}

type ResourceObj struct {
	ResourceId    string `bson:"rid"     json:"rid"`
	ResourceType  string `bson:"rtype"   json:"rtype"`
	ResourceOwner string `bson:"rowner"  json:"rowner"`
}
