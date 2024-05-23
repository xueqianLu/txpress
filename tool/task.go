package tool

import (
	"sync"
	"time"
)

type TaskHandle func (interface{})()

type Tasks struct {
	tasknum int
	handler TaskHandle
	taskpool chan interface{}
	wg sync.WaitGroup
	rate int // speed per second
}

func NewTasks(routine int, handle TaskHandle, taskpool chan interface{}) *Tasks {
	return &Tasks{
		tasknum: routine,
		handler: handle,
		taskpool: taskpool,
		rate: 0,
	}
}

func NewTasksWithSpeed(routine int, handle TaskHandle, taskpool chan interface{}, speed int) *Tasks {
	return &Tasks{
		tasknum: routine,
		handler: handle,
		taskpool: taskpool,
		rate: speed, // call handler
	}
}


func (t *Tasks) Run() {

	for i := 0; i < t.tasknum; i++ {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()

			if t.rate > 0 {
				delta := time.Nanosecond * time.Second/time.Duration(t.tasknum) / time.Duration(t.rate) // nano
				//log.Info("task delta is ", delta.Microseconds())
				t1 := time.Now()
				index := int64(0)
				for {
					select {
					case task,ok := <- t.taskpool:
						if !ok {
							return
						}

						nowPosition := time.Now().Sub(t1).Nanoseconds()
						needPosition := delta.Nanoseconds() * index
						if nowPosition < needPosition {
							//log.Info("goto wait position", "index=", index, "nowis", nowPosition, "need", needPosition)
							time.Sleep(time.Duration(needPosition - nowPosition))
						}
						index ++
						t.handler(task)
					}
				}
			} else {
				for {
					select {
					case task,ok := <- t.taskpool:
						if !ok {
							return
						}

						t.handler(task)
					}
				}
			}
		}()
	}
}

func (t *Tasks) Done() {
	t.wg.Wait()
}
