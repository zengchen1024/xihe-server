package gitlab

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

type FileCreateResult struct {
	File   string `json:"file_path"`
	Branch string `json:"branch"`
}
