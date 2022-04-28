/*
 * @Author: Tperam
 * @Date: 2022-04-28 23:43:03
 * @LastEditTime: 2022-04-29 00:13:14
 * @LastEditors: Tperam
 * @Description:
 * @FilePath: \multilock_example\lock.go
 */

package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tperam/multilock/locker"
)

var lockMessagePool sync.Pool = sync.Pool{}

type RedisLock struct {
	redis    *redis.Client
	lockname string
	lm       *lockMessage
}

// no copy
// 每次用完清除
type lockMessage struct {
	ctx context.Context
	num int
}

func (rl *RedisLock) Lock() (err error) {
	var t bool // false
	for !t {
		t, err = rl.redis.SetNX(rl.lm.ctx, rl.lockname, rl.lm.num, 1000*time.Millisecond).Result()
		if err != nil {
			rl.lm = nil
			lockMessagePool.Put(rl.lm)
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
	rl.lm = nil
	lockMessagePool.Put(rl.lm)

	return err

}

type GenerateRedisLock struct {
	redis *redis.Client
}

func (grl *GenerateRedisLock) New(lockname string) (locker.Locker, error) {
	// 尝试获取
	messageInter := lockMessagePool.Get()
	v, ok := messageInter.(*lockMessage)
	if !ok {
		v = &lockMessage{}
	}
	v.ctx = nil
	v.num = rand.Int()
	var result RedisLock
	result.redis = grl.redis
	result.lockname = lockname
	result.lm = v

	return &result, nil
}
