/*
 * @Author: Tperam
 * @Date: 2022-04-28 23:40:58
 * @LastEditTime: 2022-05-06 22:27:17
 * @LastEditors: Tperam
 * @Description:
 * @FilePath: \multilock_example\lock_test.go
 */
package multilock_redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/tperam/multilock_redis"
)

func TestRedisLock(t *testing.T) {
	ctx := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.50:30000",
		Password: "1929564872",
		DB:       0,
	})

	cmd := rdb.Get(ctx, "hh")

	fmt.Println("cmd", cmd)
	r, err := cmd.Result()
	fmt.Println("cmd result,r,err", r, err)

	m := multilock_redis.NewRedisMultilock(rdb)

	data := []string{"hi1", "hi2", "hi3", "hi4"}
	m.Do(func() (interface{}, error) {
		return businessFunc(data)
	}, data...)
}

func businessFunc(insertData []string) (interface{}, error) {
	// 假装插入数据库
	for i := range insertData {
		fmt.Println("insert data", insertData[i])
	}
	return nil, nil
}
