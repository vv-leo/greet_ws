package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisUtils struct {
}

// 存储结构体
func (RedisUtils) SetStruct(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}

	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Set(key, p, expiration).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}

// 获取结构体
func (RedisUtils) GetStruct(c *redis.Client, key string, dest interface{}) error {
	p, err := c.Get(key).Result()

	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}

func (RedisUtils) SetStructWithGob(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	return c.Set(key, buf.Bytes(), expiration).Err()
}

func (RedisUtils) GetStructWithGob(c *redis.Client, key string, dest interface{}) error {
	val, err := c.Get(key).Bytes()
	if err == redis.Nil {
		return nil
	}

	dec := gob.NewDecoder(bytes.NewReader(val))
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}
	return nil
}

// 根据规则删除缓存
func (RedisUtils) RemoveCacheByPattern(c *redis.Client, cacheKey string) (int, error) {
	var cnt = 0
	iter := c.Scan(0, cacheKey, 0).Iterator()
	for iter.Next() {
		c.Del(iter.Val())
		cnt++
	}

	return cnt, iter.Err()
}

// 判断 key 是否存在
func (RedisUtils) KeyExists(c *redis.Client, key string) (bool, error) {
	// 使用 Exists 命令检查 key 是否存在
	exists, err := c.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("检查key是否存在失败: %v", err)
	}

	// exists 返回的是 int64 类型，大于 0 表示 key 存在
	return exists > 0, nil
}

// Set 方法设置键值对，可以设置过期时间
func (RedisUtils) Set(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	// 使用 Set 命令设置键值对
	err := c.Set(key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置key失败: %v", err)
	}
	return nil
}

// 在 utils/RedisUtils.go 中添加
func (RedisUtils) Get(c *redis.Client, key string) (string, error) {
	val, err := c.Get(key).Result()
	if err == redis.Nil {
		return "", nil // key不存在
	}
	if err != nil {
		return "", fmt.Errorf("获取key失败: %v", err)
	}
	return val, nil
}
