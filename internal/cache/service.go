package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	items "github.com/broswen/kvs/internal/item"
	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	ttl int
	rdb *redis.Client
}

func New() (*CacheService, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	ttl := os.Getenv("TTL")
	ttlInSeconds, err := strconv.Atoi(ttl)
	if err != nil {
		log.Printf("unable to parse ttl: %v\n", ttl)
		ttlInSeconds = 3600
	}

	return &CacheService{
		ttl: ttlInSeconds,
		rdb: rdb,
	}, nil
}

func (cs *CacheService) IsNil() bool {
	return cs == nil
}

func (cs CacheService) Get(key string) (items.Item, error) {
	val, err := cs.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return items.Item{}, items.ErrItemNotFound{Err: err}
	}
	return items.Item{Key: key, Value: val}, nil
}

func (cs CacheService) Set(item items.Item) error {
	_, err := cs.rdb.Set(context.Background(), item.Key, item.Value, time.Duration(cs.ttl*int(time.Second))).Result()
	if err != nil {
		return items.ErrService{Err: err}
	}
	return nil
}

func (cs CacheService) Delete(key string) error {
	_, err := cs.rdb.Del(context.Background(), key).Result()
	if err != nil {
		return items.ErrItemNotFound{Err: err}
	}
	return nil
}
