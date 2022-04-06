package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gozknight.com/online-judge/internal/model"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "mmc", time.Second*10)
}

func TestRedisGet(t *testing.T) {
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}

func TestRedisGetByModel(t *testing.T) {
	val, err := model.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}
