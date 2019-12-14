package data

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type GoCacheClient struct {
	cacheClient *cache.Cache
}

func NewGoCacheClient() *GoCacheClient {
	return &GoCacheClient{
		cacheClient: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (client *GoCacheClient) Get(key string) (value interface{}, exists bool) {
	return client.cacheClient.Get(key)
}

func (client *GoCacheClient) Set(key string, value interface{}) {
	client.SetExpiration(key, value, 2*7*24*time.Hour)
}

func (client *GoCacheClient) SetExpiration(key string, value interface{}, expiration time.Duration) {
	client.cacheClient.Set(key, value, expiration)
}
