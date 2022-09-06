package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type MysqlConfig struct {
	DB_READ_HOST      string
	DB_READ_PORT      string
	DB_READ_USERNAME  string
	DB_READ_PASSWORD  string
	DB_WRITE_HOST     string
	DB_WRITE_PORT     string
	DB_WRITE_USERNAME string
	DB_WRITE_PASSWORD string
	DB_DATABASE       string
	DB_CHARSET        string
	DB_POOL_MAX       int
	DB_POOL_MIN       int
}

func LoadMysqlConfig() *MysqlConfig {
	s, err := newMysqlConfig()
	if err != nil {
		panic("mysql配置初始化失败:" + err.Error())
	}
	return s
}

func newMysqlConfig() (*MysqlConfig, error) {
	DB_POOL_MAX, _ := strconv.Atoi(os.Getenv("DB_POOL_MAX"))
	DB_POOL_MIN, _ := strconv.Atoi(os.Getenv("DB_POOL_MIN"))

	return &MysqlConfig{
		os.Getenv("DB_READ_HOST"),
		os.Getenv("DB_READ_PORT"),
		os.Getenv("DB_READ_USERNAME"),
		os.Getenv("DB_READ_PASSWORD"),
		os.Getenv("DB_WRITE_HOST"),
		os.Getenv("DB_WRITE_PORT"),
		os.Getenv("DB_WRITE_USERNAME"),
		os.Getenv("DB_WRITE_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_CHARSET"),
		DB_POOL_MAX,
		DB_POOL_MIN,
	}, nil
}
