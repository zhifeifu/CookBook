package limit

import (
	"sync"
	"time"
)

//插槽
type slot struct {
	timestamp time.Time //这个插槽的时间起点
	count     int       //这个插槽内的请求数
}

//滑动窗口
type slidingWindow struct {
	lock           sync.Mutex    //锁
	SlotDuration   time.Duration //插槽的时间长度
	WindowDuration time.Duration //滑动窗口的时间长度
	slotNum        int           //滑动窗口中的插槽个数
	windows        []*slot       //当前滑动窗口中的插槽
	maxCount       int           //滑动窗口的最大请求数
}

//构造函数
func newSlidingWindow(slotDuration time.Duration, windowDuration time.Duration, maxCount int) *slidingWindow {
	return &slidingWindow{
		SlotDuration:   slotDuration,
		WindowDuration: windowDuration,
		slotNum:        int(windowDuration / slotDuration),
		maxCount:       maxCount,
	}
}

//校验方法
func (s *slidingWindow) validate() bool {
	//加锁 解锁
	s.lock.Lock()
	defer s.lock.Unlock()
	//1.校验数量是否超限
	now := time.Now()
	timeoutOffset := -1
	for i, window := range s.windows {
		//如果插槽的时间+窗口的时间 大于当前时间
		if window.timestamp.Add(s.WindowDuration).After(now) {
			break
		}
		timeoutOffset = i
	}
	//删除窗口外的插槽
	if timeoutOffset > -1 {
		s.windows = s.windows[timeoutOffset+1:]
	}
	var check bool
	if s.windowCount() < s.maxCount {
		check = true
	}
	//2.记录本次的请求
	var lastSlot *slot
	if len(s.windows) > 0 {
		lastSlot = s.windows[len(s.windows)-1]
		// 当前时间 大于 插槽的时间 加上 插槽的跨度 小于
		if lastSlot.timestamp.Add(s.SlotDuration).Before(now) {
			lastSlot = &slot{
				timestamp: time.Now(),
				count:     1,
			}
			s.windows = append(s.windows, lastSlot)
		} else {
			lastSlot.count++
		}
	} else {
		lastSlot = &slot{
			timestamp: time.Now(),
			count:     1,
		}
		s.windows = append(s.windows, lastSlot)
	}
	return check
}

//统计滑动窗口中的总的请求数
func (s *slidingWindow) windowCount() int {
	var num int
	for _, window := range s.windows {
		num += window.count
	}
	return num
}
