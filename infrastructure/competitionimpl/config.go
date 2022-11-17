package competitionimpl

type Config struct {
	OBS OBSConfig `json:"obs" required:"true"`
}

type OBSConfig struct {
	Prefix    string `json:"prefix"`
	Bucket    string `json:"bucket"         required:"true"`
	Endpoint  string `json:"endpoint"       required:"true"`
	AccessKey string `json:"access_key"     required:"true"`
	SecretKey string `json:"secret_key"     required:"true"`
}
