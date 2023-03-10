package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance   *Config
	lock       sync.Mutex
	configPath string
)

type Config struct {
	CacheDir                string `yaml:"cacheDir"`
	Port                    int    `yaml:"port"`
	MysqlUser               string `yaml:"mysqlUser"`
	MysqlPass               string `yaml:"mysqlPass"`
	MysqlAddr               string `yaml:"mysqlAddr"`
	MysqlPort               string `yaml:"mysqlPort"`
	MysqlDBName             string `yaml:"mysqlDBName"`
	RedisAddr               string `yaml:"redisAddr"`
	RedisPass               string `yaml:"redisPass"`
	ImageOptimizationFormat string `yaml:"imageFormat"`
}

func Instance() *Config {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = load()
		}
	}
	return instance
}

func load() *Config {
	configFile, err := os.ReadFile("/etc/fuu/Percevalfile")
	if err != nil {
		overrideWithArgs()
		configFile, err = os.ReadFile(configPath)
		if err != nil {
			log.Fatalln("cannot find config file or config file not supplied")
		}
	}

	config := &Config{}

	if err == nil {
		yaml.Unmarshal(configFile, &config)
		return config
	}

	// If nothing is provided fallback to Env variables values
	fallbackToEnv(config)
	return config
}

func GetPath() string {
	return configPath
}

func fallbackToEnv(config *Config) {
	config.MysqlUser = os.Getenv("MYSQL_USER")
	config.MysqlPass = os.Getenv("MYSQL_PASS")
	config.MysqlAddr = os.Getenv("MYSQL_ADDR")
	config.MysqlPort = os.Getenv("MYSQL_PORT")
	config.MysqlDBName = os.Getenv("MYSQL_DB_NAME")
}

func overrideWithArgs() {
	flag.StringVar(&configPath, "c", "./Percevalfile", "Config file path")
	flag.Parse()
}
