package redis

// HSet set hash to redis.
func HSet(key string, field string, val string) error {
	err := GetClient().HSet(key, field, val).Err()
	if err != nil {
		return err
	}
	return nil
}

// HGet get hash from redis.
func HGet(key string, field string) (string, error) {
	val, err := GetClient().HGet(key, field).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// HGetAll get hash from redis.
func HGetAll(key string) (map[string]string, error) {
	val, err := GetClient().HGetAll(key).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return val, nil
}
