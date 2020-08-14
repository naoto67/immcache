package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/naoto67/imcache/immemcache"
	"github.com/naoto67/imcache/imredis"
	"github.com/stretchr/testify/assert"
)

var Client ICache
var ClientSet = []ICache{immemcache.NewMemcache("tcp", "localhost:11211"), imredis.NewRedis("unix", "/var/run/redis/redis-server.sock")}

type sample struct {
	Name string
}

func TestMain(m *testing.M) {
	var code int
	for _, c := range ClientSet {
		Client = c
		Client.MultiSet(map[string][]byte{"key1": []byte("\"name\""), "key2": []byte("\"name\"")})
		code = m.Run()
		Client.Flush()
	}
	os.Exit(code)
}

func Test_SingleSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		item := sample{Name: "string"}
		b, _ := json.Marshal(item)
		err := Client.SingleSet("key", b)
		assert.Nil(t, err)
	})

	t.Run("get nil key", func(t *testing.T) {
		_, err := Client.SingleGet("key100")
		assert.NotNil(t, err)
	})
}

func Test_SingleGet(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		res, err := Client.SingleGet("key1")
		assert.Nil(t, err)
		var s string
		err = json.Unmarshal(res, &s)
		assert.Nil(t, err)
		assert.Equal(t, s, "name")
	})
}

func Test_MultiSet(t *testing.T) {
	t.Run("mset", func(t *testing.T) {
		err := Client.MultiSet(map[string][]byte{"key1": []byte("\"name\""), "key2": []byte("\"name\"")})
		assert.Nil(t, err)
	})
}

func Test_MultiGet(t *testing.T) {
	t.Run("mget", func(t *testing.T) {
		_, err := Client.MultiGet([]string{"key1", "key2"})
		assert.Nil(t, err)
	})

	t.Run("mget include nil", func(t *testing.T) {
		_, err := Client.MultiGet([]string{"key1", "key2", "key3"})
		assert.Nil(t, err)
	})
}

func Test_Delete(t *testing.T) {
	t.Run("del", func(t *testing.T) {
		err := Client.SingleDelete("key1")
		assert.Nil(t, err)
	})

	t.Run("del key nil", func(t *testing.T) {
		err := Client.SingleDelete("key1000")
		assert.Nil(t, err)
	})
}

func Test_MultiDelete(t *testing.T) {
	t.Run("mdel", func(t *testing.T) {
		err := Client.MultiDelete([]string{"key1", "key2"})
		assert.Nil(t, err)
	})

	t.Run("mdel include nil", func(t *testing.T) {
		err := Client.MultiDelete([]string{"key1", "key2", "key3"})
		assert.Nil(t, err)
	})
}
