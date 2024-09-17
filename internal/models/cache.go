// Created by Anh Cao on 27.08.2024.

package models

import (
	"sync"
	"time"

	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"go.uber.org/zap"
)

type Cache struct {
	Data map[string]CacheValue
	lock sync.Mutex
}

type CacheValue struct {
	Value      interface{}
	Expiration time.Time
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a duration representing the expiration time of the value.
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAfterTimePeriod(key string, value interface{}, duration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expirationTime := time.Now().Add(duration)
	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expirationTime,
	}
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a time slot (by hour) representing the expiration time of the value
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAtTime(key string, value interface{}, expiredTime time.Time) {
	logger.Logger.Debug("set expired time for cache", zap.Time("expired-time-utc", expiredTime))
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expiredTime,
	}
}

// a method is used to retrieve a value from the cache by using a key
// It first acquires a lock on the mutex to ensure thread safety.
// Then checks if the cache contains a value for the given key and if that value has not expired.
// If the value is still valid, it returns the value and a boolean value of `true` to indicate that a valid value was found.
// If the value is not valid (means not yet cached), it returns `nil` and a boolean value of `false`.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, isValid := c.Data[key]
	if !isValid || time.Now().After(value.Expiration) {
		delete(c.Data, key)
		logger.Logger.Debug("cache was expired or not yet cached", zap.String("cache-key", key))
		return nil, false
	}
	logger.Logger.Debug("cache living time.",
		zap.Any("expired-time-utc", value.Expiration),
		zap.Time("current-time-utc", time.Now().UTC()),
	)
	return value.Value, true
}
