package protocols

import "time"

type LimiterRepo interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Increment(key string) error
}
