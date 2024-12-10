package conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Server  Server
	Elastic Elastic
	Mysql   Mysql
}

type Server struct {
	Host string
	Port string
}

type Elastic struct {
	Host         string
	Username     string
	Password     string
	MaxIdleConns int
	Timeout      time.Duration
}

type Mysql struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

func GetConfig() *Config {
	v := viper.New()
	//basePath, _ := filepath.Abs(".")
	//v.AddConfigPath(basePath)
	//v.SetConfigName(".env")
	////viper.AutomaticEnv()

	cwd, _ := os.Getwd()
	v.SetConfigFile(cwd + "/.env")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	return &Config{
		Server: Server{
			Host: v.GetString("SERVER_HOST"),
			Port: v.GetString("SERVER_PORT"),
		},
		Elastic: Elastic{
			Host:         v.GetString("ES_HOST"),
			Username:     v.GetString("ES_USERNAME"),
			Password:     v.GetString("ES_PASS"),
			MaxIdleConns: v.GetInt("MAX_IDLE_CONN"),
			Timeout:      time.Duration(v.GetInt("TIMEOUT")),
		},
		Mysql: Mysql{
			User: v.GetString("DB_USER"),
			Pass: v.GetString("DB_PASS"),
			Host: v.GetString("DB_HOST"),
			Port: v.GetString("DB_PORT"),
			Name: v.GetString("DB_NAME"),
		},
	}

}
