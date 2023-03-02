package pkg

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance *Config
	lock     = &sync.Mutex{}
)

type Config struct {
	ServerSecret            string `yaml:"serverSecret"`
	Masterpass              string `yaml:"masterpass"`
	CacheDir                string `yaml:"cacheDir"`
	WorkingDir              string `yaml:"workingDir"`
	ForceRegeneration       bool   `yaml:"regenerateThumbnailsOnBoot"`
	ThumbnailHeight         int    `yaml:"thumbnailHeight"`
	ThumbnailQuality        int    `yaml:"thumbnailQuality"`
	Port                    int    `yaml:"port"`
	UseMySQL                bool   `yaml:"useMySQL"`
	MysqlUser               string `yaml:"mysqlUser"`
	MysqlPass               string `yaml:"mysqlPass"`
	MysqlAddr               string `yaml:"mysqlAddr"`
	MysqlPort               string `yaml:"mysqlPort"`
	MysqlDBName             string `yaml:"mysqlDBName"`
	RedisAddr               string `yaml:"redisAddr"`
	RedisPass               string `yaml:"redisPass"`
	ImageOptimizationFormat string `yaml:"imageOptimizationFormat"`
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
			log.Println("Cannot find config file, fallbacking to ENVvariables")
		}
	}

	config := Config{}

	if err == nil {
		yaml.Unmarshal(configFile, &config)
		return &config
	}

	// If nothing is provided fallback to Env variables values
	fallbackToEnv(&config)
	overrideWithArgs(&config)
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

	height, err := strconv.Atoi(os.Getenv("THUMBNAIL_HEIGHT"))
	if err != nil {
		height = 450
	}
	config.ThumbnailHeight = height

	quality, err := strconv.Atoi(os.Getenv("THUMBNAIL_QUALITY"))
	if err != nil {
		quality = 75
	}
	config.ThumbnailQuality = quality
}

func overrideWithArgs(config *Config) {
	flag.StringVar(&config.Masterpass, "M", "adminadmin", "Main user password")
	flag.StringVar(&config.ServerSecret, "S", "secret", "Signing secret")
	flag.StringVar(&config.WorkingDir, "w", ".", "Pictures directory")

	flag.IntVar(&config.ThumbnailHeight, "th", 450, "Thumbnails height (px)")
	flag.IntVar(&config.ThumbnailQuality, "tq", 75, "Thumbnails quality (0-100]")
	flag.IntVar(&config.Port, "p", 4456, "Where server will listen at")

	flag.Parse()
}
