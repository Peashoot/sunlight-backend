package utils

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/peashoot/sunlight/config"
)

func RedisSet(key string, value interface{}) error {
	return config.RedisDB.Set(key, value, 0).Err()
}

func RedisSetT(key string, value interface{}, expiration time.Duration) error {
	return config.RedisDB.Set(key, value, expiration).Err()
}

func RedisRemove(key ...string) error {
	return config.RedisDB.Del(key...).Err()
}

func RedisGet[T any](key string) (T, error) {
	var ret T
	val, err := config.RedisDB.Get(key).Result()
	if err != nil {
		return ret, err
	}
	ret, err = jsonUnmall[T](val)
	return ret, err
}

func jsonUnmall[T any](value interface{}) (T, error) {
	var ret T
	var ok bool
	if reflect.TypeOf(ret).Kind() == reflect.String {
		if ret, ok = value.(T); ok {
			return ret, nil
		}
	}
	var str string
	if str, ok = value.(string); ok {
		err := json.Unmarshal([]byte(str), &ret)
		return ret, err
	}
	return ret, errors.New("type of mismatch")
}

func RedisTryLock(key string, expiration time.Duration) (bool, error) {
	result := config.RedisDB.SetNX(key, time.Now(), expiration)
	return result.Result()
}

func RedisTryUnlock(key ...string) (bool, error) {
	result := config.RedisDB.Del(key...)
	return result.Val() > 0, result.Err()
}

func RedisExpire(key string, expiration time.Duration) error {
	return config.RedisDB.Expire(key, expiration).Err()
}

func RedisHSet(key, mapKey string, value interface{}) error {
	return config.RedisDB.HSet(key, mapKey, value).Err()
}

func RedisHGet[T any](key, mapKey string) (T, error) {
	var ret T
	val, err := config.RedisDB.HGet(key, mapKey).Result()
	if err != nil {
		return ret, err
	}
	ret, err = jsonUnmall[T](val)
	return ret, err
}

func RedisHDel(key string, mapKeys ...string) error {
	return config.RedisDB.HDel(key, mapKeys...).Err()
}

func RedisExists(keys ...string) (bool, error) {
	result := config.RedisDB.Exists(keys...)
	return result.Val() > 0, result.Err()
}

func RedisScan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return config.RedisDB.Scan(cursor, match, count).Result()
}
