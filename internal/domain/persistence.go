package domain

type Persistence interface {
	Set(key string, value string, ttl int) error
	Get(key string) (string, error)
}
