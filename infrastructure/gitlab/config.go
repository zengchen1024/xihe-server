package gitlab

type Config struct {
	OBS             OBSConfig `json:"obs"              required:"true"`
	LFSPath         string    `json:"lfs_path"         required:"true"`
	Endpoint        string    `json:"endpoint"         required:"true"`
	RootToken       string    `json:"root_token"       required:"true"`
	GraphqlEndpoint string    `json:"graphql_endpoint" required:"true"`
	DefaultBranch   string    `json:"default_branch"`
	DownloadExpiry  int       `json:"download_expiry"`
}

func (cfg *Config) SetDefault() {
	cfg.DefaultBranch = "main"
	cfg.DownloadExpiry = 3600
}

type OBSConfig struct {
	Bucket    string `json:"bucket"     required:"true"`
	Endpoint  string `json:"endpoint"   required:"true"`
	AccessKey string `json:"access_key" required:"true"`
	SecretKey string `json:"secret_key" required:"true"`
}
