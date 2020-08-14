package main

type ICache interface {
	SingleSet(key string, value []byte) error
	MultiSet(set map[string][]byte) error
	SingleGet(key string) ([]byte, error)
	MultiGet(keys []string) ([][]byte, error)
	SingleDelete(key string) error
	MultiDelete(keys []string) error
	Flush() error
}
