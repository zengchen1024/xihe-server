package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId             = "id"
	fieldPId            = "pid"
	fieldBio            = "bio"
	fieldJob            = "job"
	fieldDesc           = "desc"
	fieldLevel          = "level"
	fieldDetail         = "detail"
	fieldCoverId        = "cover_id"
	fieldCommit         = "commit"
	fieldExpiry         = "expiry"
	fieldRepoId         = "repo_id"
	fieldTags           = "tags"
	fieldKinds          = "kinds"
	fieldStatus         = "status"
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
	fieldProjects       = "projects"
	fieldRId            = "rid"
	fieldTId            = "tid"
	fieldROwner         = "rowner"
	fieldRType          = "rtype"
	fieldUpdatedAt      = "updated_at"
	fieldDownloadCount  = "download_count"
	fieldFirstLetter    = "fl"
	fieldPhase          = "phase"
	fieldTeams          = "teams"
	fieldRepos          = "repos"
	fieldOrder          = "order"
	fieldEnabled        = "enabled"
	fieldCompetitors    = "competitors"
	fieldSubmissions    = "submissions"
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
	Level    int    `bson:"level"      json:"level"`
	Name     string `bson:"name"       json:"name"`
	FL       byte   `bson:"fl"         json:"fl"`
	Desc     string `bson:"desc"       json:"desc"`
	CoverId  string `bson:"cover_id"   json:"cover_id"`
	RepoType string `bson:"repo_type"  json:"repo_type"`
	// set omitempty to avoid set it to null occasionally.
	Tags     []string `bson:"tags"       json:"tags,omitempty"`
	TagKinds []string `bson:"kinds"      json:"kinds,omitempty"`
}

type dModel struct {
	Owner string      `bson:"owner" json:"owner"`
	Items []modelItem `bson:"items" json:"-"`
}

type modelItem struct {
	Id        string `bson:"id"         json:"id"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	ModelPropertyItem `bson:",inline"`

	// RelatedDatasets is not allowd to be set,
	// So, don't marshal it to avoid setting it occasionally.
	RelatedDatasets []ResourceIndex `bson:"datasets" json:"-"`
	RelatedProjects []ResourceIndex `bson:"projects" json:"-"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"           json:"-"`
	LikeCount     int `bson:"like_count"        json:"-"`
	DownloadCount int `bson:"download_count"    json:"-"`
}

type ModelPropertyItem struct {
	FL       byte     `bson:"fl"         json:"fl"`
	Level    int      `bson:"level"      json:"level"`
	Name     string   `bson:"name"       json:"name"`
	Desc     string   `bson:"desc"       json:"desc"`
	RepoType string   `bson:"repo_type"  json:"repo_type"`
	Tags     []string `bson:"tags"       json:"tags,omitempty"`
	TagKinds []string `bson:"kinds"      json:"kinds,omitempty"`
}

type dDataset struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []datasetItem `bson:"items" json:"-"`
}

type datasetItem struct {
	Id        string `bson:"id"         json:"id"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	DatasetPropertyItem `bson:",inline"`

	RelatedModels   []ResourceIndex `bson:"models"   json:"-"`
	RelatedProjects []ResourceIndex `bson:"projects" json:"-"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"               json:"-"`
	LikeCount     int `bson:"like_count"            json:"-"`
	DownloadCount int `bson:"download_count"    json:"-"`
}

type DatasetPropertyItem struct {
	FL       byte     `bson:"fl"         json:"fl"`
	Level    int      `bson:"level"      json:"level"`
	Name     string   `bson:"name"       json:"name"`
	Desc     string   `bson:"desc"       json:"desc"`
	RepoType string   `bson:"repo_type"  json:"repo_type"`
	Tags     []string `bson:"tags"       json:"tags,omitempty"`
	TagKinds []string `bson:"kinds"      json:"kinds,omitempty"`
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
	Items []dDomainTags `bson:"items"    json:"items"`
}

type dDomainTags struct {
	Name   string  `bson:"name"          json:"name"`
	Domain string  `bson:"domain"        json:"domain"`
	Tags   []dTags `bson:"tags"          json:"tags"`
}

type dTags struct {
	Kind string   `bson:"kind"           json:"kind"`
	Tags []string `bson:"tags"           json:"tags"`
}

type dTraining struct {
	Owner         string `bson:"owner"   json:"owner"`
	ProjectId     string `bson:"pid"     json:"pid"`
	ProjectName   string `bson:"name"    json:"name"`
	ProjectRepoId string `bson:"rid"     json:"rid"`
	Version       int    `bson:"version" json:"-"`

	Items []trainingItem `bson:"items"   json:"-"`
}

type trainingItem struct {
	Id             string      `bson:"id"            json:"id"`
	Name           string      `bson:"name"          json:"name"`
	Desc           string      `bson:"desc"          json:"desc"`
	CodeDir        string      `bson:"code_dir"      json:"code_dir"`
	BootFile       string      `bson:"boot_file"     json:"boot_file"`
	Compute        dCompute    `bson:"compute"       json:"compute"`
	Inputs         []dInput    `bson:"inputs"        json:"inputs"`
	EnableAim      bool        `bson:"aim"           json:"aim"`
	EnableOutput   bool        `bson:"output"        json:"output"`
	Env            []dKeyValue `bson:"env"           json:"env"`
	Hypeparameters []dKeyValue `bson:"parameters"    json:"parameters"`
	CreatedAt      int64       `bson:"created_at"    json:"created_at"`

	Job       dJobInfo   `bson:"job"     json:"job"`
	JobDetail dJobDetail `bson:"detail"  json:"detail"`
}

type dCompute struct {
	Type    string `bson:"type"          json:"type"`
	Flavor  string `bson:"flavor"        json:"flavor"`
	Version string `bson:"version"       json:"version"`
}

type dKeyValue struct {
	Key   string `bson:"key"             json:"key"`
	Value string `bson:"value"           json:"value"`
}

type dInput struct {
	Key    string `bson:"key"            json:"key"`
	Type   string `bson:"type"           json:"type"`
	User   string `bson:"user"           json:"user"`
	File   string `bson:"file"           json:"file"`
	RepoId string `bson:"rid"            json:"rid"`
}

type dJobInfo struct {
	Endpoint  string `bson:"endpoint"    json:"endpoint"`
	JobId     string `bson:"job_id"      json:"job_id"`
	LogDir    string `bson:"log"         json:"log"`
	AimDir    string `bson:"aim"         json:"aim"`
	OutputDir string `bson:"output"      json:"output"`
}

type dJobDetail struct {
	Duration   int    `bson:"duration"   json:"duration,omitempty"`
	Status     string `bson:"status"     json:"status,omitempty"`
	LogPath    string `bson:"log"        json:"log,omitempty"`
	AimPath    string `bson:"aim"        json:"aim,omitempty"`
	OutputPath string `bson:"output"     json:"output,omitempty"`
}

type dInference struct {
	Owner       string `bson:"owner"   json:"owner"`
	ProjectId   string `bson:"pid"     json:"pid"`
	ProjectName string `bson:"name"    json:"name"`
	LastCommit  string `bson:"commit"  json:"commit"`
	Version     int    `bson:"version" json:"-"`

	Items []inferenceItem `bson:"items"  json:"-"`
}

type inferenceItem struct {
	Id        string `bson:"id"          json:"id,omitempty"`
	Expiry    int64  `bson:"expiry"      json:"expiry,omitempty"`
	Error     string `bson:"error"       json:"error,omitempty"`
	AccessURL string `bson:"url"         json:"url,omitempty"`
}

type dEvaluate struct {
	Owner      string `bson:"owner"       json:"owner"`
	ProjectId  string `bson:"pid"         json:"pid"`
	TrainingId string `bson:"tid"         json:"tid"`
	Version    int    `bson:"version"     json:"-"`

	Items []evaluateItem `bson:"items"  json:"-"`
}

type evaluateItem struct {
	// must set omitempty, otherwise it will be override to empty
	Id                string   `bson:"id"          json:"id,omitempty"`
	Type              string   `bson:"type"        json:"type,omitempty"`
	MomentumScope     []string `bson:"momentum"    json:"momentum,omitempty"`
	BatchSizeScope    []string `bson:"bsize"       json:"bsize,omitempty"`
	LearningRateScope []string `bson:"rate"        json:"rate,omitempty"`
	Expiry            int64    `bson:"expiry"      json:"expiry,omitempty"`
	Error             string   `bson:"error"       json:"error,omitempty"`
	AccessURL         string   `bson:"url"         json:"url,omitempty"`
}

type DCompetition struct {
	Id         string `bson:"id"              json:"id"`
	Name       string `bson:"name"            json:"name"`
	Desc       string `bson:"desc"            json:"desc"`
	Host       string `bson:"host"            json:"host"`
	Phase      string `bson:"phase"           json:"phase"`
	Bonus      int    `bson:"bonus"           json:"bonus"`
	Status     string `bson:"status"          json:"status"`
	Duration   string `bson:"duration"        json:"duration"`
	Doc        string `bson:"doc"             json:"doc"`
	Forum      string `bson:"forum"             json:"forum"`
	Poster     string `bson:"poster"          json:"poster"`
	DatasetDoc string `bson:"dataset_doc"     json:"dataset_doc"`
	DatasetURL string `bson:"dataset_url"     json:"dataset_url"`
	SmallerOk  bool   `bson:"order"           json:"order"`
	Enabled    bool   `bson:"enabled"         json:"enabled"`

	Teams       []dTeam            `bson:"teams"       json:"-"`
	Repos       []dCompetitionRepo `bson:"repos"       json:"-"`
	Competitors []dCompetitor      `bson:"competitors" json:"-"`
	Submissions []dSubmission      `bson:"submissions" json:"-"`
}

type dCompetitor struct {
	Name     string            `bson:"name"      json:"name"`
	City     string            `bson:"city"      json:"city"`
	Email    string            `bson:"email"     json:"email"`
	Phone    string            `bson:"phone"     json:"phone"`
	Account  string            `bson:"account"   json:"account"`
	Identity string            `bson:"identity"  json:"identity"`
	Province string            `bson:"province"  json:"province"`
	Detail   map[string]string `bson:"detail"    json:"detail"`
	TeamId   string            `bson:"tid"       json:"tid"`
	TeamRole string            `bson:"role"      json:"role"`
}

type dTeam struct {
	Id   string `bson:"id"      json:"id"`
	Name string `bson:"name"    json:"name"`
}

type dSubmission struct {
	Id         string  `bson:"id"          json:"id"`
	TeamId     string  `bson:"tid"         json:"tid"`     // if it is submitted by team, set it.
	Individual string  `bson:"account"     json:"account"` // if it is submitted by individual, set it.
	Status     string  `bson:"status"      json:"status"`
	OBSPath    string  `bson:"path"        json:"path"`
	SubmitAt   int64   `bson:"submit_at"   json:"submit_at"`
	Score      float32 `bson:"score"       json:"score"`
}

type dCompetitionRepo struct {
	TeamId     string `bson:"tid"         json:"tid"`     // if it is submitted by team, set it.
	Individual string `bson:"account"     json:"account"` // if it is submitted by individual, set it.
	Owner      string `bson:"owner"       json:"owner"`
	Repo       string `bson:"repo"        json:"repo"`
}

type dLuoJia struct {
	Owner string       `bson:"owner" json:"owner"`
	Items []luojiaItem `bson:"items" json:"-"`
}

type luojiaItem struct {
	Id        string `bson:"id"         json:"id"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
}
