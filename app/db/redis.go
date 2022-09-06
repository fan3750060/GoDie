package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"goframe/app/logger"
	"goframe/config"
	"time"
)

var redisConfig *config.RedisConfig
var redisClient *redis.Pool

type Redis interface {
	Set(key string, value string, args ...int) bool
	Get(key string) string
	Del(key string) bool
	GetConn() redis.Conn
}

type RedisClass struct {
}

func init() {
	redisConfig = config.LoadRedisConfig()

	dns := fmt.Sprintf("%s:%s", redisConfig.REDIS_HOST, redisConfig.REDIS_PORT)

	redisClient = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dns)
			if err != nil {
				logger.Logger.Println(fmt.Sprintf("Failed to connect redis , err : %s \n", err.Error()))
				return nil, err
			}
			c.Do("SELECT", redisConfig.REDIS_DB)

			if redisConfig.REDIS_AUTH != "" {
				if _, err := c.Do("AUTH", redisConfig.REDIS_AUTH); err != nil {
					logger.Logger.Println(fmt.Sprintf("Failed to auth redis , err : %s \n", err.Error()))
					return c, nil
				}
			}

			logger.Logger.Println("SUCCESS: redis连接成功")

			return c, nil
		},

		//DialContext:     nil,

		//TestOnBorrow:    nil,

		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: redisConfig.REDIS_POOL_MIN,

		//最大的激活连接数，表示同时最多有N个连接
		MaxActive: redisConfig.REDIS_POOL_MAX,

		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 180 * time.Second,

		//Wait:            false,

		//MaxConnLifetime: 0,
	}

	logger.Logger.Println("SUCCESS: redis连接成功")

}

/**
 * [设置缓存 返回布尔值]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (redisClass RedisClass) Set(key string, value string, args ...int) bool {

	redisConn := redisClient.Get()
	defer redisConn.Close()

	timeOut := 0
	if len(args) != 0 {
		if args[0] <= 0 {
			timeOut = 0
		} else {
			timeOut = args[0]
		}
	}

	var err error
	var res interface{}

	if timeOut == 0 {
		res, err = redisConn.Do("SET", key, value)
	} else {
		res, err = redisConn.Do("SETEX", key, timeOut, value)
	}

	if err == nil {
		if res == "OK" {
			logger.Logger.Println(fmt.Sprintf("redis-succes SET: %s EXPIRE : %d \n", key, timeOut))
			return true
		}

		logger.Logger.Println(fmt.Sprintf("redis-error SET: %s \n", key))
		return false
	} else {
		logger.Logger.Println(fmt.Sprintf("redis-error:%s \n", err.Error()))
		return false
	}
}

/**
 * [获取缓存 返回字符串]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (redisClass RedisClass) Get(key string) string {
	redisConn := redisClient.Get()
	defer redisConn.Close()

	res, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("redis获取key失败，error：%s \n", err.Error()))
		return ""
	}

	return res
}

/**
 * [删除缓存 返回布尔值]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (redisClass RedisClass) Del(key string) bool {
	redisConn := redisClient.Get()
	defer redisConn.Close()

	res, err := redisConn.Do("DEL", key)

	if err != nil {
		logger.Logger.Println(fmt.Sprintf("redis删除key错误，error：%s \n", err.Error()))
		return false
	}

	if res == int64(1) {
		logger.Logger.Println(fmt.Sprintf("redis删除key成功，key：%s \n", key))
		return true
	}

	logger.Logger.Println(fmt.Sprintf("redis删除key失败，key：%s \n", key))
	return false
}

/**
 * [redis连接 返回客户端连接]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (redisClass RedisClass) GetConn() redis.Conn {
	redisConn := redisClient.Get()
	defer redisConn.Close()
	return redisConn
}

//...其他操作自行实现 用到再说
