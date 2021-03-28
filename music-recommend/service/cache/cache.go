package cache

type Cache interface {
	Set(key string, val interface{}, expire int, priority Priority)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Priority int

const (
	One   Priority = 1
	Two   Priority = 2
	Three Priority = 3
	Four  Priority = 4
	Five  Priority = 5
	Six   Priority = 6
	Seven Priority = 7
	Eight Priority = 8
	Nine  Priority = 9
	Ten   Priority = 10
	zero  Priority = 0

	NoExpire   = -1
	zeroExpire = 0

	Hour = 60 * 60
)
