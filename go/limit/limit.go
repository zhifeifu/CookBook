package limit

import (
	"sync"
	"time"
)

//限流器-计数器
type Counter struct {
	count     int           //当前请求数
	limit     int           //限制值
	beginTime time.Time     //开始时间
	cycle     time.Duration //周期
	lock      sync.Mutex    //锁
}

func NewCounter(limit int, cycle time.Duration) *Counter {
	return &Counter{
		count:     0,
		limit:     limit,
		beginTime: time.Now(),
		cycle:     cycle,
	}
}

//检查
func (l *Counter) check() bool {
	//加锁
	l.lock.Lock()
	//解锁
	defer l.lock.Unlock()
	//当前请求数是否大于限制值
	if l.count >= l.limit {
		//当前时间
		now := time.Now()
		if now.Sub(l.beginTime) >= l.cycle {
			l.beginTime = now
			l.count = 0
			return true
		}
		return false
	} else {
		l.count++
		return true
	}
}
