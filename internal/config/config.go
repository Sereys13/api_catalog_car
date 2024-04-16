package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Listen struct {
	Type   string
	BindIp string
	Port   string
}

type Storage struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Sslmode  string
}

type Logs struct{
	PathLog string
	NameFileLOg string
}


type Config struct {
	Listen  Listen
	Storage Storage
	Logs Logs
	UrlApiCarInfo string
}

func GetConfig() (*Config, error) {
	err := godotenv.Load("../config.env")
	if err != nil {
		return nil, err
	}
	return &Config{Listen: Listen{
		Type:   os.Getenv("LISTEN_TYPE"),
		BindIp: os.Getenv("LISTEN_BIND_IP"),
		Port:   os.Getenv("LISTEN_PORT")},
		Storage: Storage{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			Sslmode:  os.Getenv("DB_SSLMODE")},
		Logs: Logs{
			PathLog: os.Getenv("PATCH_LOG"),
			NameFileLOg: os.Getenv("NAME_LOG")},
		UrlApiCarInfo: os.Getenv("URL_API_CAR_INFO"),}, nil
}