package config

import (
	"fmt"
	"io/ioutil"

	"github.com/scriptscat/cloudcat/internal/pkg/database"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/cache"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mode     string
	Database *database.Config
	KvDB     *kvdb.Config `yaml:"kvDb"`
	Cache    *cache.Config
	OAuth    struct {
		BBS OAuth `yaml:"bbs"`
	} `yaml:"oauth"`
	Addr string `yaml:"addr"`
}

type OAuth struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

func Init(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config read error: %v", err)

	}
	ret := &Config{}
	err = yaml.Unmarshal(file, ret)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return ret, nil
}
