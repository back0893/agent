package main

import (
	"agent/src"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ti := src.NewTimingWheel(context.Background())
	var wg sync.WaitGroup

	wg.Add(4)

	ml := make([]int64, 0)
	var id int64
	id = ti.AddTimer(time.Now(), time.Second*1, func() {
		t.Logf("定时器执行1")
	})
	ml = append(ml, id)

	ti.AddTimer(time.Now(), time.Second*2, func() {
		t.Logf("定时器执行2")
	})
	ml = append(ml, id)

	ti.AddTimer(time.Now(), time.Second*3, func() {
		t.Logf("定时器执行3")
	})
	ml = append(ml, id)

	ti.AddTimer(time.Now().Add(5*time.Second), 0, func() {
		for _, id := range ml {
			ti.Cancel(id)
			wg.Done()
		}
	})

	go ti.Start()
	go func() {
		defer wg.Done()
		notify := make(chan os.Signal)
		signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case <-notify:
				{
					t.Log("中断")
					return
				}
			}
		}

	}()
	wg.Wait()

}
