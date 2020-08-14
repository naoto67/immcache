package main

type ICache interface {
	SingleSet(key string, value []byte) error
	SingleGet(key string) ([]byte, error)
	MultiSet(set map[string][]byte) error
	MultiGet(keys []string) ([][]byte, error)
	Flush() error
}
