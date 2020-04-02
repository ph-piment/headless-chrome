package redis

// LPush left most push list to redis.
func LPush(key string, val []string) error {
	err := GetClient().LPush(key, val).Err()
	if err != nil {
		return err
	}
	return nil
}

// RPush right most push list to redis.
func RPush(key string, val []string) error {
	err := GetClient().RPush(key, val).Err()
	if err != nil {
		return err
	}
	return nil
}

// LSet set list to redis.
func LSet(key string, i int64, val string) error {
	err := GetClient().LSet(key, i, val).Err()
	if err != nil {
		return err
	}
	return nil
}

// LRange get list from redis.
func LRange(key string, start, stop int64) ([]string, error) {
	val, err := GetClient().LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}
	return val, nil
}

// AllRange get all list from redis.
func AllRange(key string) ([]string, error) {
	val, err := LRange(key, 0, -1)
	if err != nil {
		return []string{}, err
	}
	return val, nil
}

// LIndex get list from redis.
func LIndex(key string, i int64) (string, error) {
	val, err := GetClient().LIndex(key, i).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
