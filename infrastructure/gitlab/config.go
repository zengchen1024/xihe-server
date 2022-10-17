package gitlab

type Config struct {
	OBS            OBSConfig `json:"obs"              required:"true"`
	Endpoint       string    `json:"endpoint"         required:"true"`
	RootToken      string    `json:"root_token"       required:"true"`
	DefaultBranch  string    `json:"default_branch"   required:"true"`
	LFSPath        string    `json:"lfs_path"         required:"true"`
	DownloadExpiry int       `json:"download_expiry"`
}

type OBSConfig struct {
	Endpoint  string `json:"endpoint"   required:"true"`
	AccessKey string `json:"access_key" required:"true"`
	SecretKey string `json:"secret_key" required:"true"`
	Bucket    string `json:"bucket"     required:"true"`
}
