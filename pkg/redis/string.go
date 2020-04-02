package redis

import (
	"time"
)

// SetString set string to redis.
func SetString(key string, val string) error {
	err := SetStringWithExpire(key, val, 0)
	if err != nil {
		return err
	}
	return nil
}

// SetStringWithExpire set string to redis with expire.
func SetStringWithExpire(key string, val string, expire time.Duration) error {
	err := GetClient().Set(key, val, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetString get string from redis.
func GetString(key string) (string, error) {
	val, err := GetClient().Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// DelString delete string from redis.
func DelString(key string) error {
	err := GetClient().Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}
