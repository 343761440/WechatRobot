package common

import (
	"errors"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/redis.v5"
)

var (
	once      sync.Once
	gInstance *redis.Client
)

func initRedisClient() (*redis.Client, error) {
	//初始化redis连接
	redisAddress := viper.GetString("REDIS.ADDRESS")
	redisPort := viper.GetString("REDIS.PORT")
	if len(redisAddress) == 0 || len(redisPort) == 0 {
		return nil, errors.New("redis address is empty")
	}
	rds := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{redisAddress + ":" + redisPort},
	})
	if ping, err := rds.Ping().Result(); err != nil {
		panic("redis " + ping + ", err " + err.Error())
	}
	return rds, nil
}

func GetRedisClient() (*redis.Client, error) {
	var rerr error
	once.Do(func() {
		rdb, err := initRedisClient()
		if err != nil {
			rerr = err
			return
		}
		gInstance = rdb
	})
	return gInstance, rerr
}
