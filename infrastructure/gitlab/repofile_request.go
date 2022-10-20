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
	Trees graphqlNodes `json:"trees"`
}

type graphqlNodes struct {
	Nodes []graphqlNode `json:"nodes"`
}

type graphqlNode struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	LFSOid string `json:"lfsOid"`
}

func (d *graphqlData) toRepoPathItems() (r []platform.RepoPathItem) {
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
