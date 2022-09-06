package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type RedisConfig struct {
	REDIS_CLUSTER bool
	REDIS_HOST    string
	REDIS_AUTH    string
	REDIS_PORT    string
	REDIS_DB      int
	REDIS_POOL_MAX int
	REDIS_POOL_MIN int
}

func LoadRedisConfig() *RedisConfig {
	s, err := newRedisConfig()
	if err != nil {
		panic("redis配置初始化失败:" + err.Error())
	}
	return s
}

func newRedisConfig() (*RedisConfig, error) {
	cluster, _ := strconv.Atoi(os.Getenv("REDIS_CLUSTER"))

	REDIS_CLUSTER := cluster != 0

	REDIS_DB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	REDIS_POOL_MAX, _ := strconv.Atoi(os.Getenv("REDIS_POOL_MAX"))

	REDIS_POOL_MIN, _ := strconv.Atoi(os.Getenv("REDIS_POOL_MIN"))

	return &RedisConfig{
		REDIS_CLUSTER,
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_AUTH"),
		os.Getenv("REDIS_PORT"),
		REDIS_DB,
		REDIS_POOL_MAX,
		REDIS_POOL_MIN,
	}, nil
}
