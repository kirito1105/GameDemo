package myRedis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type MyRedis struct {
	rdb *redis.Client
}

const (
	ROOM_EXP_TIME = time.Second * 60
)

var myRedis *MyRedis
var once sync.Once

func GetMyRedisCon() *MyRedis {
	once.Do(func() {
		myRedis = &MyRedis{rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		})}
	})
	return myRedis
}

func (r *MyRedis) SetEX(key string, value string) error {
	r.rdb.Set(context.Background(), key, value, ROOM_EXP_TIME)
	return nil
}

func (r *MyRedis) HSetEX(key string, field map[string]interface{}) error {
	txn := r.rdb.TxPipeline()
	txn.HSet(context.Background(), key, field)
	txn.Expire(context.Background(), key, ROOM_EXP_TIME)
	txn.Exec(context.Background())
	return nil
}
