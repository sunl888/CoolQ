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
