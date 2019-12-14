package data

import "time"

type CacheProvider interface {
	Get(key string) (value interface{}, exists bool)
	Set(key string, value interface{})
	SetExpiration(key string, value interface{}, exp time.Duration)
}
