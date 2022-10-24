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

	Content  string `json:"content"`
	Encoding string `json:"encoding,omitempty"`
}

type graphqlResult struct {
	Data graphqlData `json:"data"`
}

type graphqlData struct {
	Project graphqlProject `json:"project"`
}

type graphqlProject struct {
	Repo graphqlRepo `json:"repository"`
}

type graphqlRepo struct {
	Tree  graphqlTree  `json:"tree"`
	Blobs graphqlBlobs `json:"blobs"`
}

type graphqlTree struct {
	Blobs      graphqlBlobs  `json:"blobs"`
	Trees      graphqlTrees  `json:"trees"`
	LastCommit graphqlCommit `json:"lastCommit"`
}

type graphqlBlobs struct {
	Nodes []graphqlNode `json:"nodes"`
}

type graphqlTrees struct {
	Nodes []graphqlNode `json:"nodes"`
}

type graphqlCommit struct {
	SHA string `json:"sha"`
}

type graphqlNode struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	LFSOid string `json:"lfsOid"`
}

func (d *graphqlResult) toRepoPathItems() (r []platform.RepoPathItem) {
	files := d.Data.Project.Repo.Tree.Blobs.Nodes
	dirs := d.Data.Project.Repo.Tree.Trees.Nodes

	total := len(files) + len(dirs)
	if total == 0 {
		return
	}

	r = make([]platform.RepoPathItem, total)

	for i := range files {
		item := &files[i]

		r[i] = platform.RepoPathItem{
			Path:      item.Path,
			Name:      item.Name,
			IsLFSFile: item.LFSOid != "",
		}
	}

	if len(dirs) == 0 {
		return
	}

	v := r[len(files):]

	for i := range dirs {
		item := &dirs[i]

		v[i] = platform.RepoPathItem{
			Path:  item.Path,
			Name:  item.Name,
			IsDir: true,
		}
	}

	return
}
