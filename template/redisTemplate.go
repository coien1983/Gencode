package template

var RedisTemplate = `package sysinit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

func RedisInit(setting *AppConfig) (err error) {

	redisConf := setting.RedisC

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       redisConf.Db,
		PoolSize: redisConf.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return

}

func GetRedis() *redis.Client {
	return rdb
}

func RedisClose() {
	_ = rdb.Close()
}

//计算Redis的剩余时间
func RedisTtl(key string) (time.Duration, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.TTL(keyR).Result()
}

func RedisIsExists(key string) (int64, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.Exists(keyR).Result()
}

// RedisReadString 读取字符串
func RedisReadString(key string) (string, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	return rdb.Get(keyR).Result()
}

// RedisWriteString 写入字符串
func RedisWriteString(key string, value interface{}, expiredTime int64) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	newTime := time.Duration(expiredTime) * time.Second
	return rdb.Set(keyR, value, newTime).Err()
}

// RedisReadStruct 读取结构体
func RedisReadStruct(key string, obj interface{}) error {

	if data, err := RedisReadString(key); err == nil {
		return json.Unmarshal([]byte(data), obj)
	} else {
		return err
	}
}

// RedisWriteStruct 写入结构体
func RedisWriteStruct(key string, obj interface{}, expiredTime int64) error {
	data, err := json.Marshal(obj)
	if err == nil {
		return RedisWriteString(key, string(data), expiredTime)
	} else {
		return err
	}
}

// RedisDelete 删除键
func RedisDelete(key string) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.Del(keyR).Err()
}

// RedisPopQueue 从队列中获取数据
func RedisPopQueue(key string) (number string, err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	val, err := rdb.RPop(keyR).Result()
	if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

// RedisPushQueue 推送数据入队列
func RedisPushQueue(key string, value interface{}) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.LPush(keyR, value).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("入列失败")
	}

	return nil
}

func RedisIncr(key string) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.Incr(keyR).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisDecr(key string) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.Decr(keyR).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisHSet(key string, field, data string, dayTime time.Time) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	err := rdb.HSet(keyR, field, data).Val()
	if !err {
		return errors.New("操作失败")
	}

	_, _ = rdb.ExpireAt(keyR, dayTime).Result()

	return nil
}

func RedisHGet(key, data string) (string, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	data, err := rdb.HGet(keyR, data).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func RedisHIncrBy(key string, data string, expiredTime int64) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.HIncrBy(keyR, data, expiredTime).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisExpiredKey(key string, t time.Time) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	_, _ = rdb.ExpireAt(keyR, t).Result()
}
`
