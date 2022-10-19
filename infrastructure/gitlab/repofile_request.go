package gitlab

import "github.com/opensourceways/xihe-server/domain/platform"

type CommitInfo struct {
	Branch      string `json:"branch"           required:"true"`
	Message     string `json:"commit_message"   required:"true"`
	AuthorName  string `json:"author_name"      required:"false"`
	AuthorEmail string `json:"author_email"     required:"false"`
}

type FileCreateOption struct {
	CommitInfo

	Content string `json:"content"          required:"true"`
}

type graphqlData struct {
	Data graphqlProject `json:"data"`
}

type graphqlProject struct {
	Project graphqlRepo `json:"project"`
}

type graphqlRepo struct {
	Repo graphqlTree `json:"repository"`
}

type graphqlTree struct {
	Tree graphqlBlobs `json:"tree"`
}

type graphqlBlobs struct {
	Blobs graphqlNodes `json:"blobs"`
}

type graphqlNodes struct {
	Nodes []graphqlNode `json:"nodes"`
}

type graphqlNode struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	LFSOid string `json:"lfsOid"`
}

func (d *graphqlData) toRepoPathItems() (r []platform.RepoPathItem) {
	v := d.Data.Project.Repo.Tree.Blobs.Nodes

	r = make([]platform.RepoPathItem, len(v))

	for i := range v {
		r[i] = platform.RepoPathItem{
			Name:      v[i].Name,
			IsDir:     v[i].Type == "tree",
			IsLFSFile: v[i].LFSOid != "",
		}
	}

	return
}
