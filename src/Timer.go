package src

import (
	"context"
	"time"
)

/**
定时器,定时检查心跳状态
*/

type Timer struct {
	ctx    context.Context
	timers []context.CancelFunc
}

func (t *Timer) Add(fn func() bool, interval time.Duration) int {
	id := len(t.timers) + 1
	timeInterval := time.NewTicker(interval)
	ctx, cancelFunc := context.WithCancel(t.ctx)
	go func(ctx context.Context) {
		select {
		case <-timeInterval.C:
			stop := fn()
			if stop == true {
				t.Cancel(id)
			}
		case <-t.ctx.Done():
			return
		}
	}(ctx)
	t.timers = append(t.timers, cancelFunc)
	return id
}

func (t *Timer) Cancel(id int) {
	if id > len(t.timers) {
		return
	}
	fn := t.timers[id]
	fn()
}

func NewTimer(ctx context.Context) *Timer {
	return &Timer{
		ctx:    ctx,
		timers: make([]context.CancelFunc, 0),
	}
}
