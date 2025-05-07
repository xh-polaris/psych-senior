package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"

	"github.com/zeromicro/go-zero/core/service"

	"github.com/zeromicro/go-zero/core/conf"
)

var config *Config

type SMTP struct {
	Username string
	Password string
	Host     string
	Port     int
	Alert    string
}
type Config struct {
	service.ServiceConf
	ListenOn string
	State    string
	Auth     Auth
	Mongo    struct {
		URL string
		DB  string
	}
	Cache               cache.CacheConf
	Redis               *redis.RedisConf
	RabbitMQ            RabbitMQ
	SMTP                SMTP
	BaiLianChat         BaiLianChat
	BaiLianShanghaiChat BaiLianShanghaiChat
	BaiLianReport       BaiLianReport
	VolcTts             VolcTts
	VolcAsr             VolcAsr
	VolcNoModelTts      VolcNoModelTts
}

type Auth struct {
	SecretKey    string
	PublicKey    string
	AccessExpire int64
}

type RabbitMQ struct {
	Url string
}

type BaiLianChat struct {
	AppId  string
	ApiKey string
}

type BaiLianShanghaiChat struct {
	AppId  string
	ApiKey string
}

type BaiLianReport struct {
	AppId  string
	ApiKey string
}

type VolcTts struct {
	Url        string
	AppKey     string
	AccessKey  string
	Speaker    string
	ResourceId string
}

type VolcNoModelTts struct {
	Url       string
	AppKey    string
	AccessKey string
	Speaker   string
	Cluster   string
	Lang      string
}

type VolcAsr struct {
	Url        string
	AppKey     string
	AccessKey  string
	ResourceId string
}

func NewConfig() (*Config, error) {
	c := new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, c)
	if err != nil {
		return nil, err
	}
	err = c.SetUp()
	if err != nil {
		return nil, err
	}
	config = c
	return c, nil
}

func GetConfig() *Config {
	return config
}
