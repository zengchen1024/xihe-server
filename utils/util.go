package utils

import (
	"io/ioutil"
	"os"
	"time"

	"sigs.k8s.io/yaml"
)

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := []byte(os.ExpandEnv(string(b)))

	return yaml.Unmarshal(content, cfg)
}

func Now() int64 {
	return time.Now().Unix()
}

func ToDate(n int64) string {
	return time.Unix(n, 0).Format("2006-01-02")
}

func Expiry(expiry int64) int64 {
	return time.Now().Add(time.Second * time.Duration(expiry)).Unix()
}
