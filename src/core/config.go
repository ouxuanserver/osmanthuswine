package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Port          string `json:"port"`
	Host          string `json:"host"`
	CrossDomain   string `json:"cross_domain"`
	ApiRouter     string `json:"api_router"`
	PostMaxMemory int64  `json:"post_max_memory"`
	Db            struct {
		Host        string            `json:"host"`
		Port        string            `json:"port"`
		User        string            `json:"user"`
		Password    string            `json:"password"`
		Name        string            `json:"name"`
		Prefix      string            `json:"prefix"`
		MaxOpenConn int               `json:"max_open_conn"`
		Params      map[string]string `json:"params"`
		Debug       bool              `json:"debug"`
	} `json:"db"`
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
	UpdateDir  string `json:"update_dir"`
	UpdatePath string `json:"update_path"`
}

var instanceConfig *Config
var configFile = "./config.json"
var privateConfigFile = "./private.json"

func SetConfigFile(c string) {
	configFile = c
}

func SetConfig(c *Config) {
	instanceConfig = c
}

func GetInstanceConfig() *Config {
	if instanceConfig == nil {
		instanceConfig = &Config{
			Host:          "localhost",
			Port:          "8808",
			ApiRouter:     "/Api/*",
			CrossDomain:   "*",
			PostMaxMemory: 1024 * 1024 * 10,
			Db: struct {
				Host        string            `json:"host"`
				Port        string            `json:"port"`
				User        string            `json:"user"`
				Password    string            `json:"password"`
				Name        string            `json:"name"`
				Prefix      string            `json:"prefix"`
				MaxOpenConn int               `json:"max_open_conn"`
				Params      map[string]string `json:"params"`
				Debug       bool              `json:"debug"`
			}{
				Host:        "localhost",
				Port:        "3306",
				User:        "root",
				Password:    "root",
				Name:        "test",
				Prefix:      "",
				MaxOpenConn: 500,
				Params: map[string]string{
					"charset":   "utf8mb4",
					"parseTime": "true",
				},
				Debug: true,
			},
			Redis: struct {
				Addr     string `json:"addr"`
				Password string `json:"password"`
				Db       int    `json:"db"`
			}{
				Addr:     "localhost:6379",
				Password: "",
				Db:       0,
			},
			UpdateDir:  "",
			UpdatePath: os.Args[0] + "_update",
		}

		instanceConfig.ReadConfig(configFile)
		instanceConfig.ReadPrivateConfig(privateConfigFile)
	}
	return instanceConfig
}

func (c *Config) ReadConfig(file string) {
	configText, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("配置文件错误,启动失败:", err.Error())
		os.Exit(0)
	}
	err = json.Unmarshal(configText, c)
	if err != nil {
		log.Println("配置文件错误,启动失败:", err.Error())
		os.Exit(0)
	}
}

func (c *Config) ReadPrivateConfig(file string) {
	configText, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("未加载", privateConfigFile, ":", err.Error())
		return
	}
	err = json.Unmarshal(configText, c)
	if err != nil {
		log.Println("未加载", privateConfigFile, ":", err.Error())
	}
}
