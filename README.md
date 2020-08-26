## 使い方

基本的に、main以下に置いているファイルをwgetして使ってください。

```sh
wget https://raw.githubusercontent.com/naoto67/immcache/master/main/cache.go
```

### Redisを使用する場合

```go
var cacheClient *redisClient
cacheClient = NewRedis("tcp", "localhost:6379")

data := []byte("1")
cacheClient.SingleSet("key", data)
```

### 上記をMemcachedに入れ替える

```go
// var cacheClient *redisClient
// cacheClient = NewRedis("tcp", "localhost:6379")
var cacheClient *memcacheClient
cacheClient = NewMemcache("tcp", "localhost:6379")

data := []byte("1")
cacheClient.SingleSet("key", data)
```

### cacheClientのインターフェース
redisClientとmemcacheClientはそれぞれFetchConnというレスポンスの型が違うものが存在するため、一旦このInterfaceを使用していない。
```
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
}
```

### 注意点

MultiGetをする時、一部nilの要素が含まれる場合があるため、nilはスキップしてパースする必要がある
```
res, _ := MultiGet([]string{"aaa", "bbb"})
for i := range res {
  if res[i] == nil {
    continue
  }
  json.Marshal...
}
```
