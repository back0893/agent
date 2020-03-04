package src

import (
	"agent/src/agent/model"
	"sync"
)

type TaskQueue struct {
	queue []*model.Service
	cond  *sync.Cond
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue: make([]*model.Service, 0),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}
func (task *TaskQueue) Push(action *model.Service) {
	task.cond.L.Lock()
	defer task.cond.L.Unlock()
	task.queue = append(task.queue, action)
	task.cond.Signal() //唤醒一个wait的协成
}

/**
等待一直有值返回,如果请注意
*/
func (task *TaskQueue) Pop() *model.Service {
	task.cond.L.Lock()
	defer task.cond.L.Unlock()
	var n int
	for {
		n = len(task.queue)
		if n > 0 {
			break
		}
		task.cond.Wait()
	}
	action := task.queue[n-1]
	task.queue = task.queue[:n-1]
	return action
}
