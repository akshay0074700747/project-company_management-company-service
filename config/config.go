package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBhost     string
	DBname     string
	DBport     string
	DBuser     string
	DBpassword string
	EndPoint   string
	AccessKey  string
	SecretKey  string
}

func LoadConfigurations() (Config, error) {

	if err := godotenv.Load(".env"); err != nil {
		return Config{}, err
	}

	var conf Config

	conf.DBhost = os.Getenv("dbhost")
	conf.DBport = os.Getenv("dbport")
	conf.DBname = os.Getenv("dbname")
	conf.DBpassword = os.Getenv("dbpassword")
	conf.DBuser = os.Getenv("dbuser")
	conf.EndPoint = os.Getenv("endpoint")
	conf.SecretKey = os.Getenv("secretKey")
	conf.AccessKey = os.Getenv("accessKey")

	return conf, nil
}
