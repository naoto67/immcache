package main

type ICache interface {
	SingleSet(key string, value []byte) error
	SingleSetNX(key string, value []byte) (int, error)
	MultiSet(set map[string][]byte) error
	SingleGet(key string) ([]byte, error)
	MultiGet(keys []string) ([][]byte, error)
	SingleDelete(key string) error
	MultiDelete(keys []string) error
	Increment(key string, delta uint64) (int, error)
	Flush() error

	FetchConn() interface{}
}
