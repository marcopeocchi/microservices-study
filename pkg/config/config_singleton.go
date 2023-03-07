package pkg

import (
	"flag"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance   *Config
	lock       = &sync.Mutex{}
	configPath string
)

type Config struct {
	ServerSecret            string   `yaml:"serverSecret"`
	Masterpass              string   `yaml:"masterpass"`
	CacheDir                string   `yaml:"cacheDir"`
	WorkingDir              string   `yaml:"workingDir"`
	ForceRegeneration       bool     `yaml:"regenerateThumbnailsOnBoot"`
	ThumbnailHeight         int      `yaml:"thumbnailHeight"`
	ThumbnailQuality        int      `yaml:"thumbnailQuality"`
	Port                    int      `yaml:"port"`
	UseMySQL                bool     `yaml:"useMySQL"`
	MysqlUser               string   `yaml:"mysqlUser"`
	MysqlPass               string   `yaml:"mysqlPass"`
	MysqlAddr               string   `yaml:"mysqlAddr"`
	MysqlPort               string   `yaml:"mysqlPort"`
	MysqlDBName             string   `yaml:"mysqlDBName"`
	RedisAddr               string   `yaml:"redisAddr"`
	RedisPass               string   `yaml:"redisPass"`
	ImageOptimizationFormat string   `yaml:"imageOptimizationFormat"`
	ImageProcessors         []string `yaml:"imageProcessors"`
	TLSCertPath             string   `yaml:"tlsCertPath"`
	JaegerEndpoint          string   `yaml:"jaegerEndpoint"`
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
	configFile, err := os.ReadFile("./Fuufile")

	if err != nil {
		configFile, err = os.ReadFile("/etc/fuu/Fuufile")
		if err != nil {
			configFile, err = os.ReadFile(configPath)
			if err != nil {
				log.Fatalln("cannot find config file or config file not supplied")
			}
		}
	}

	config := Config{}

	if err == nil {
		yaml.Unmarshal(configFile, &config)
		return &config
	}

	// If nothing is provided fallback to Env variables values
	fallbackToEnv(&config)
	overrideWithArgs()
	return &config
}

func fallbackToEnv(config *Config) {
	config.Masterpass = os.Getenv("MASTERPASS")
	config.ServerSecret = os.Getenv("SECRET")
	config.WorkingDir = os.Getenv("WORKDIR")
	config.MysqlUser = os.Getenv("MYSQL_USER")
	config.MysqlPass = os.Getenv("MYSQL_PASS")
	config.MysqlAddr = os.Getenv("MYSQL_ADDR")
	config.MysqlPort = os.Getenv("MYSQL_PORT")
	config.MysqlDBName = os.Getenv("MYSQL_DB_NAME")
}

func overrideWithArgs() {
	flag.StringVar(&configPath, "c", "./Fuufile", "Config file path")
	flag.Parse()
}
