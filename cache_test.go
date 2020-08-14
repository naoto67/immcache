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
var ClientSet = []ICache{
	immemcache.NewMemcache("tcp", "localhost:11211"),
	imredis.NewRedis("unix", "/var/run/redis/redis-server.sock"),
}

type TestCase struct {
	name string
	call func(t *testing.T)
}

type sample struct {
	Name string
}

func TestMain(m *testing.M) {
	var code int
	code = m.Run()
	os.Exit(code)
}

func TestAll(t *testing.T) {
	testCases := []TestCase{
		TestCase{
			name: "SingleSet",
			call: func(t *testing.T) {
				item := sample{Name: "string"}
				b, _ := json.Marshal(item)
				err := Client.SingleSet("key", b)
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "SingleSetNX",
			call: func(t *testing.T) {
				b := []byte("\"sam\"")
				stored, err := Client.SingleSetNX("key", b)
				assert.Nil(t, err)
				assert.Equal(t, stored, 1)
			},
		},
		TestCase{
			name: "SingleSetNX: exists key",
			call: func(t *testing.T) {
				b := []byte("\"sam\"")
				stored, err := Client.SingleSetNX("key1", b)
				assert.Nil(t, err)
				assert.Equal(t, stored, 0)
			},
		},
		TestCase{
			name: "SingleGet",
			call: func(t *testing.T) {
				res, err := Client.SingleGet("key1")
				assert.Nil(t, err)
				var s string
				json.Unmarshal(res, &s)
				assert.Equal(t, s, "name")
			},
		},
		TestCase{
			name: "MultiSet",
			call: func(t *testing.T) {
				err := Client.MultiSet(map[string][]byte{"key1": []byte("\"name\""), "key2": []byte("\"name\"")})
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "MultiGet",
			call: func(t *testing.T) {
				_, err := Client.MultiGet([]string{"key1", "key2"})
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "MultiGet: include not exists key",
			call: func(t *testing.T) {
				_, err := Client.MultiGet([]string{"key1", "key2", "key100"})
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "SingleDelete:",
			call: func(t *testing.T) {
				err := Client.SingleDelete("key1")
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "SingleDelete: not exists key",
			call: func(t *testing.T) {
				err := Client.SingleDelete("key100")
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "MultiDelete",
			call: func(t *testing.T) {
				err := Client.MultiDelete([]string{"key1", "key2"})
				assert.Nil(t, err)
			},
		},
		TestCase{
			name: "MultiDelete: include not exists key",
			call: func(t *testing.T) {
				err := Client.MultiDelete([]string{"key1", "key2", "key100"})
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		for _, Client = range ClientSet {
			// fixtures
			Client.MultiSet(map[string][]byte{"key1": []byte("\"name\""), "key2": []byte("\"name\"")})
			t.Run(testCase.name, testCase.call)
			Client.Flush()
		}
	}
}
