package common

import (
	"time"
)

// Stopwatch 用于计时
type Stopwatch struct {
	start time.Time
}

// Start 开始计时
func (w *Stopwatch) Start() {
	w.start = time.Now()
}

// ElapsedMilliseconds 获取消耗的毫秒数
func (w *Stopwatch) ElapsedMilliseconds() int64 {
	return time.Now().Sub(w.start).Nanoseconds() / 1e6
}
