package redisclient

import (
	"time"
)

// SetKey sets a key with optional TTL
func SetKey(key string,field string, value string, ttl time.Duration) error {
		if err := RDB.HSet(Ctx, key, field, value).Err(); err != nil {
		return err
	}

	// 2️⃣ Set TTL if provided
	if ttl > 0 {
		if err := RDB.Expire(Ctx, key, ttl).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetKey gets a value by key
func GetKey(key string) (map[string]string, error) {
	return RDB.HGetAll(Ctx, key).Result()
} 
func GetTTL(key string) (time.Duration, error) {
	return RDB.TTL(Ctx, key).Result()
}
