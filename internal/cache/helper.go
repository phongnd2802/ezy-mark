package cache

import "time"

type TTLStatus int

const (
	TTLExpired    TTLStatus = -2
	TTLPersistent TTLStatus = -1
	TTLHasValue   TTLStatus = 1
)

func CheckTTL(ttl time.Duration) TTLStatus {
	if ttl == -2 {
		return TTLExpired
	} else if ttl == -1 {
		return TTLPersistent
	}
	return TTLHasValue
}