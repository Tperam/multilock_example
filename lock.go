/*
 * @Author: Tperam
 * @Date: 2022-04-28 23:43:03
 * @LastEditTime: 2022-05-05 23:51:41
 * @LastEditors: Tperam
 * @Description:
 * @FilePath: \multilock_example\lock.go
 */

package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tperam/multilock/locker"
)

type RedisLock struct {
	redis    *redis.Client
	lockname string
	lm       lockMessage
}

type lockMessage struct {
	ctx context.Context
	num int
}

func (rl *RedisLock) Lock() (err error) {
	var t bool // false
	for !t {
		if rl.redis == nil {
			panic("redis error")
		}
		t, err = rl.redis.SetNX(rl.lm.ctx, rl.lockname, rl.lm.num, 1000*time.Millisecond).Result()
		if err != nil {
			return err
		}
		// 缓冲
		time.Sleep(time.Millisecond)
	}
	return err
}

func (rl *RedisLock) Unlock() error {
	//TODO 判断并执行删除，通过lua脚本
	_, err := rl.redis.Del(context.TODO(), rl.lockname).Result()
	// 可能还需要做处理，比如判断是因为什么原因导致的
	return err

}

type GenerateRedisLock struct {
	redis *redis.Client
}

func NewGenerateRedisLock(redis *redis.Client) *GenerateRedisLock {
	return &GenerateRedisLock{
		redis: redis,
	}
}

func (grl *GenerateRedisLock) New(lockname string) (locker.Locker, error) {
	v := &RedisLock{}
	// 尝试获取
	lm := lockMessage{}
	lm.ctx = context.TODO()
	lm.num = rand.Int()

	v.redis = grl.redis
	v.lockname = lockname
	v.lm = lm

	return v, nil
}
