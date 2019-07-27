package config

import (
	"encoding/json"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	MessageHandlerUrl string `json:"message_handler_url"`
	Token             string `json:"token"`
	NotifyUrl         string `json:"notify_url"` // "notify_url": "https://oapi.dingtalk.com/robot/send?access_token=a36cf190b4021ff97964396f23d2986a83566fdd608bfd59d1d1979459c5c4b7",
}

func LoadConfig() (*Config, error) {
	file, err := getFile("config.json")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

//获取插件文件的路径，若路径不存在则顺便创建
func getFile(name string) (string, error) {
	appDir := cqp.GetAppDir()
	err := os.MkdirAll(appDir, os.ModeDir)
	return filepath.Join(appDir, name), err
}
