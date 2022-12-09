package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 微信公众号APP_ID
	WxAppId string `json:"wx_app_id"`
	// 微信公众号Token
	WxToken string `json:"wx_token"`
	// 启动服务端口
	HttpPort string `json:"port"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		WxAppId := os.Getenv("WxAppId")
		WxToken := os.Getenv("WxToken")
		HttpPort := os.Getenv("HttpPort")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if WxAppId != "" {
			config.WxAppId = WxAppId
		}
		if WxToken != "" {
			config.WxToken = WxToken
		}
		if HttpPort != "" {
			config.HttpPort = HttpPort
		}
	})
	return config
}
