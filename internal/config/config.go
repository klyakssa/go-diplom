package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type WebServerConfig struct {
	Port int
}

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type LoggingConfiguration struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
}

type DBConfig struct {
	ConnectionString string `mapstructure:"connection-string"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
}

type Config struct {
	Debug   bool                  `mapstructure:"debug"`
	App     *AppConfig            `mapstructure:"app"`
	Logging *LoggingConfiguration `mapstructure:"logging"`
	Web     *WebServerConfig      `mapstructure:"web"`
	PostDB  *DBConfig             `mapstructure:"postdb"`
	JWT     *JWTConfig            `mapstructure:"jwt"`
}

var C *Config = new(Config)

func InitConfiguration() *Config {
	initConfig()
	return C
}

func initConfig() {
	loadDefault()
	loadFile()
	loadEnv()

	if viper.GetBool("debug") {
		viper.SetDefault("logging.level", "debug")
	}

	viper.Unmarshal(C)
}

func loadEnv() {
	viper.SetEnvPrefix("GOMART")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

const (
	AppName = "gophermart"
)

func loadDefault() {
	viper.SetDefault("debug", false)
	viper.SetDefault("app.name", AppName)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.path", "logs")
	viper.SetDefault("logging.max-size", 500)
	viper.SetDefault("logging.max-backups", 3)
	viper.SetDefault("logging.max-age", 30)

	viper.SetDefault("web.port", 8100)

	viper.SetDefault("postdb.connection-string", "postgres://test:11@localhost:5432/diplom?sslmode=disable")

	viper.SetDefault("jwt.secret", "ADJG1HAJD5GADHS3GDKAHJD2GASJHD5KA")
	viper.SetDefault("jwt.expire", "1h")
}

func loadFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", viper.GetString(AppName)))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", viper.GetString(AppName)))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/config/", viper.GetString(AppName)))

	var fileLookupError viper.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			if err := viper.WriteConfigAs("./config.json"); err != nil {
				fmt.Printf("Error writing config file: %v\n", err)
				panic(err)
			}
		} else {
			fmt.Println(fmt.Printf("Error reading config file, %s. Use default only.\n", err))
		}
	}
}
