package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jhawk7/go-searchme/internal/pkg/common"
	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisClient struct {
	svc *redis.Client
}

type GetValue struct {
	Key string
}

type KVPair struct {
	Key   string
	Value interface{}
}

type DeleteKeys struct {
	Keys []string
}

func InitClient() *RedisClient {
	var redisClient RedisClient
	redisClient.svc = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	//test redis connection
	pong, err := redisClient.svc.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to get response from redis server; [error: %v]", err))
	} else {
		log.Infof("Redis PING Response: %v", pong)
	}

	return &redisClient
}

func (redisClient *RedisClient) GetValue(ctx context.Context, key string) (value []interface{}, miss bool, err error) {
	data, getErr := redisClient.svc.Get(ctx, key).Result()
	if getErr != nil {
		if getErr.Error() != redis.Nil.Error() {
			err = fmt.Errorf("failed to retrieve value for key %v", key)
			miss = false
		} else {
			common.LogInfo(fmt.Sprintf("cache miss; [key: %v]", key))
			miss = true
		}
		return
	}

	common.LogInfo(fmt.Sprintf("cache hit, [key: %v]", key))

	if uErr := json.Unmarshal([]byte(data), &value); uErr != nil {
		err = fmt.Errorf("failed to unmarshall cache data; %v", uErr)
		return
	}

	return
}

func (redisClient *RedisClient) Store(ctx context.Context, kv KVPair) (err error) {
	bytes, mErr := json.Marshal(kv.Value)
	if mErr != nil {
		err = fmt.Errorf("failed to marshal kv value of messages; %v", mErr)
		return
	}

	if storeErr := redisClient.svc.Set(ctx, kv.Key, string(bytes), time.Hour).Err(); storeErr != nil {
		err = fmt.Errorf("unable to store kv pair; [error: %v]", storeErr)
	} else {
		common.LogInfo(fmt.Sprintf("kv pair stored [key: %v]", kv.Key))
	}
	return
}

func (redisClient *RedisClient) Delete(ctx context.Context, keys DeleteKeys) (err error) {
	if dErr := redisClient.svc.Del(ctx, keys.Keys...).Err(); dErr != nil {
		err = fmt.Errorf("unable to delete keys %v; [error: %v]", keys.Keys, dErr)
	}
	return
}
