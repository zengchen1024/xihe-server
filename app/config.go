package app

var appConfig Config

func Init(cfg *Config) {
	appConfig = *cfg
}

type Config struct {
	WuKongMaxLikeNum int `json:"wukong_max_like_num"     required:"true"`
	FinetuneMaxNum   int `json:"finetune_max_num"        required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.WuKongMaxLikeNum <= 0 {
		cfg.WuKongMaxLikeNum = 10
	}

	if cfg.FinetuneMaxNum <= 0 {
		cfg.FinetuneMaxNum = 5
	}
}
