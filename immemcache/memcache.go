package immemcache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheClient struct {
	client *memcache.Client
}

// protocol unused.
// if host contains "/", it will be used unix socket.
func NewMemcache(protocol, host string) *memcacheClient {
	client := memcache.New(fmt.Sprintf("%s", host))
	err := client.Ping()
	fmt.Println(err)

	return &memcacheClient{
		client: client,
	}
}

func (mc *memcacheClient) FetchConn() interface{} {
	return mc.client
}

func (mc *memcacheClient) SingleSet(key string, value []byte) error {
	return mc.client.Set(&memcache.Item{Key: key, Value: value})
}

func (mc *memcacheClient) SingleGet(key string) ([]byte, error) {
	item, err := mc.client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (mc *memcacheClient) MultiSet(set map[string][]byte) error {
	var err error
	for k, v := range set {
		err = mc.SingleSet(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mc *memcacheClient) MultiGet(keys []string) ([][]byte, error) {
	res := make([][]byte, len(keys))
	itemMap, err := mc.client.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	for key, _ := range itemMap {
		res = append(res, itemMap[key].Value)
	}
	return res, nil
}

func (mc *memcacheClient) SingleDelete(key string) error {
	mc.client.Delete(key)
	return nil
}

func (mc *memcacheClient) MultiDelete(keys []string) error {
	for _, key := range keys {
		mc.SingleDelete(key)
	}

	return nil
}

func (mc *memcacheClient) SingleSetNX(key string, value []byte) (int, error) {
	err := mc.client.Add(&memcache.Item{Key: key, Value: value})
	if err == memcache.ErrNotStored {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (mc *memcacheClient) Increment(key string, delta uint64) (int, error) {
	v, err := mc.client.Increment(key, delta)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (mc *memcacheClient) Flush() error {
	return mc.client.FlushAll()
}
