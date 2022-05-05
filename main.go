/*
 * @Author: Tperam
 * @Date: 2022-04-28 23:40:58
 * @LastEditTime: 2022-05-05 23:33:47
 * @LastEditors: Tperam
 * @Description:
 * @FilePath: \multilock_example\main.go
 */
package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tperam/multilock"
	"github.com/tperam/multilock/algorithm"
	"github.com/tperam/multilock/lockcore"
	"github.com/tperam/multilock/locker"
)

func main() {
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

	m := multilock.NewMultilock(algorithm.NewBubbleSort(), lockcore.NewMapLockCore(locker.NewGenerateMutex(), 10), lockcore.NewMapLockCore(NewGenerateRedisLock(rdb), 10))

	data := []string{"hi1", "hi2", "hi3", "hi4"}
	m.Do(func() (interface{}, error) {
		return BusinessFunc(data)
	}, data...)
}

func BusinessFunc(insertData []string) (interface{}, error) {
	// 假装插入数据库
	for i := range insertData {
		fmt.Println("insert data", insertData[i])
	}
	return nil, nil
}
