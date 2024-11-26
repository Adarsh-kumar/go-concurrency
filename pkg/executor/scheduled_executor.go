package executor

import (
	"context"
	"time"

	"github.com/Adarsh-kumar/go-concurrency/pkg/container"
	"github.com/Adarsh-kumar/go-concurrency/pkg/logger"
	"github.com/Adarsh-kumar/go-concurrency/pkg/task"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type ScheduledExecutor struct {
	taskQueue      *container.BoundedTaskQueue
	maxThredsSema  *semaphore.Weighted
	pushNotifyChan chan struct{}
	log            *zap.SugaredLogger
}

func NewScheduledExecutor(maxThreads int) *ScheduledExecutor {
	l := logger.GetLogger()
	l.Info("Initializing ScheduledExecutor")

	st := &ScheduledExecutor{
		taskQueue:      container.NewBoundedTaskQueue(100),
		maxThredsSema:  semaphore.NewWeighted(int64(maxThreads)),
		pushNotifyChan: make(chan struct{}),
		log:            l,
	}
	go st.startExecutor()
	return st
}

func (st *ScheduledExecutor) ScheduleTask(t task.Task) {
	st.log.Infof("Scheduling one-shot task: %s", t.Name())
	scheduledTask := &task.ScheduledTask{
		Task:           t,
		SubmissionTime: time.Now(),
		Policy:         task.OneShotTaskPolicy{},
	}
	st.taskQueue.EnqueTask(scheduledTask)
	go func() {
		st.pushNotifyChan <- struct{}{}
	}()
}

func (st *ScheduledExecutor) ScheduleTaskWithFixedDelay(t task.Task, delay time.Duration) {
	st.log.Infof("Scheduling task with fixed delay: %s, delay: %v", t.Name(), delay)
	scheduledTask := &task.ScheduledTask{
		Task:           t,
		SubmissionTime: time.Now(),
		Policy: task.DelayedTaskPolicy{
			Delay: delay,
		},
	}
	st.taskQueue.EnqueTask(scheduledTask)
	go func() {
		st.pushNotifyChan <- struct{}{}
	}()
}

func (st *ScheduledExecutor) ScheduleTaskAtFixedRate(t task.Task, interval time.Duration) {
	st.log.Infof("Scheduling task at fixed rate: %s, interval: %v", t.Name(), interval)
	scheduledTask := &task.ScheduledTask{
		Task:           t,
		SubmissionTime: time.Now(),
		Policy: task.FixedRateTaskPolicy{
			Rate: interval,
		},
	}
	st.taskQueue.EnqueTask(scheduledTask)
	go func() {
		st.pushNotifyChan <- struct{}{}
	}()
}

func (st *ScheduledExecutor) startExecutor() {
	st.log.Debug("Starting the executor")
	for {
		t := st.taskQueue.Top()
		if t == nil {
			st.log.Debug("Task queue is empty, waiting for tasks...")
			continue
		}

		if t.GetScheduledTime().Before(time.Now()) {
			st.log.Infof("Executing task: %s", t.Task.Name())
			topTask := st.taskQueue.DequeTask()
			st.processTask(topTask)
		} else {
			delay := t.GetScheduledTime().Sub(time.Now())
			st.log.Infof("Next task (%s) scheduled in %v", t.Task.Name(), delay)
			timerChan := time.After(delay)
			select {
			case <-timerChan:
				st.log.Infof("Scheduled time reached for task: %s", t.Task.Name())
				continue
			case <-st.pushNotifyChan:
				st.log.Debug("New task notification received")
				continue
			}
		}
	}
}

func (st *ScheduledExecutor) processTask(t *task.ScheduledTask) {
	ctx := context.TODO()
	st.maxThredsSema.Acquire(ctx, 1)
	go func() {
		defer st.maxThredsSema.Release(1)
		t.LastScheduledTime = time.Now()

		if _, ok := t.Policy.(task.FixedRateTaskPolicy); ok {
			st.log.Debugf("Re-enqueuing fixed-rate task: %s", t.Task.Name())
			st.taskQueue.EnqueTask(t)
			go func() {
				st.pushNotifyChan <- struct{}{}
			}()
		}

		st.log.Infof("Running task: %s", t.Task.Name())
		err := t.Task.Run()
		if err != nil {
			st.log.Errorf("Task %s failed: %v", t.Task.Name(), err)
		} else {
			st.log.Infof("Task %s completed successfully", t.Task.Name())
		}
	}()
}
