package src

import "sync"

type Action struct {
	DeviceId string   //agent编号
	Action   string   //动作
	Args     []string //参数
}

type TaskQueue struct {
	queue []*Action
	cond  *sync.Cond
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue: make([]*Action, 0),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}
func (task *TaskQueue) Push(action *Action) {
	task.cond.L.Lock()
	defer task.cond.L.Unlock()
	task.queue = append(task.queue, action)
	task.cond.Signal() //唤醒一个wait的协成
}

/**
等待一直有值返回,如果请注意
*/
func (task *TaskQueue) Pop() *Action {
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
