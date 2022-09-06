package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type MongoDbConfig struct {
	MONGODB_WRITE_HOST      string
	MONGODB_READ_HOST       string
	REPLICASETNAME          string
	MONGODB_PORT            string
	MONGODB_DATABASE        string
	MONGODB_USERNAME        string
	MONGODB_PASSWORD        string
	MONGODB_AUTHDB          string
	MONGODB_READ_PREFERENCE string
}

func LoadMongoDbConfig() *MongoDbConfig {
	s, err := newMongoDbConfig()
	if err != nil {
		panic("mongodb配置初始化失败:" + err.Error())
	}
	return s
}

func newMongoDbConfig() (*MongoDbConfig, error) {
	return &MongoDbConfig{
		os.Getenv("MONGODB_WRITE_HOST"),
		os.Getenv("MONGODB_READ_HOST"),
		os.Getenv("REPLICASETNAME"),
		os.Getenv("MONGODB_PORT"),
		os.Getenv("MONGODB_DATABASE"),
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_AUTHDB"),
		os.Getenv("MONGODB_READ_PREFERENCE"),
	}, nil
}
