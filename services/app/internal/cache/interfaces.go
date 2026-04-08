package cache

import "time"

type ICache interface {
	Get(key string) (any, bool)
	Set(key string, value any, duration time.Duration)
}
