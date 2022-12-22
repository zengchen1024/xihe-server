package app

var appConfig Config

func Init(cfg *Config) {
	appConfig = *cfg
}

type Config struct {
	WuKongMaxLikeNum int `json:"wukong_max_like_num"     required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.WuKongMaxLikeNum <= 0 {
		cfg.WuKongMaxLikeNum = 10
	}
}
