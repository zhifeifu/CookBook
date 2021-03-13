package limit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //上次向桶中放令牌的时间的时间戳，单位为秒

	lock sync.Mutex
}

//构造函数
func NewTokenBucket(rate int64, capacity int64) *TokenBucket {
	return &TokenBucket{
		rate:         rate,
		capacity:     capacity,
		tokens:       0,
		lastTokenSec: time.Now().Unix(),
	}
}

//校验
func (t *TokenBucket) validate() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	now := time.Now().Unix()
	t.tokens = t.tokens + (now-t.lastTokenSec)*t.rate
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	t.lastTokenSec = now
	if t.tokens > 0 {
		t.tokens--
		return true
	} else {
		return false
	}
}
