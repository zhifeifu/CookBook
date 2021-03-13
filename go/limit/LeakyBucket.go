package limit

import (
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	//桶的容量
	capacity float64
	//出桶的速率(每秒)
	rate float64
	//当前桶中的数量
	water float64
	//上一次漏桶的时间
	lastLeakMs int64
	//锁
	lock sync.Mutex
}

//构造函数
func NewLeakyBucket(capacity float64, rate float64) *LeakyBucket {
	return &LeakyBucket{
		water:      0,
		rate:       rate,
		capacity:   capacity,
		lastLeakMs: time.Now().UnixNano() / 1e6,
	}
}

//校验
func (l *LeakyBucket) Validate() bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	//计算桶的容量
	//当前的微妙
	nowMs := time.Now().UnixNano() / 1e6
	l.water = math.Max(0, l.water-(float64(nowMs-l.lastLeakMs)*l.rate/1000))
	l.lastLeakMs = nowMs
	if l.water+1 <= l.capacity {
		l.water++
		return true
	} else {
		return false
	}
}
