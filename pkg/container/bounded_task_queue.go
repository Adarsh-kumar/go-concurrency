package container

import (
	"container/heap"
	"sync"

	"github.com/Adarsh-kumar/go-concurrency/pkg/logger"
	"github.com/Adarsh-kumar/go-concurrency/pkg/task"
)

type BoundedTaskQueue struct {
	pq          *PriorityQueue
	size        int
	capacity    int
	lock        *sync.Mutex
	emptyWaitFn *sync.Cond
	fullWaitFn  *sync.Cond
}

func NewBoundedTaskQueue(capacity int) *BoundedTaskQueue {
	log := logger.GetLogger()
	log.Debugf("Initializing BoundedTaskQueue with capacity: %d", capacity)

	mutex := &sync.Mutex{}
	pq := &PriorityQueue{}
	heap.Init(pq)
	return &BoundedTaskQueue{
		pq:          pq,
		capacity:    capacity,
		lock:        mutex,
		emptyWaitFn: sync.NewCond(mutex),
		fullWaitFn:  sync.NewCond(mutex),
	}
}

func (bq *BoundedTaskQueue) EnqueTask(t *task.ScheduledTask) {
	log := logger.GetLogger()
	bq.lock.Lock()
	defer bq.lock.Unlock()

	for bq.size == bq.capacity {
		log.Debugf("Queue is full, waiting to enqueue task: %s", t.Task.Name())
		bq.fullWaitFn.Wait()
	}

	heap.Push(bq.pq, t)
	bq.size++
	log.Debugf("Enqueued task: %s, current size: %d", t.Task.Name(), bq.size)
	bq.emptyWaitFn.Signal()
}

func (bq *BoundedTaskQueue) DequeTask() *task.ScheduledTask {
	log := logger.GetLogger()
	bq.lock.Lock()
	defer bq.lock.Unlock()

	for bq.size == 0 {
		log.Debug("Queue is empty, waiting to dequeue task")
		bq.emptyWaitFn.Wait()
	}

	t := heap.Pop(bq.pq)
	bq.size--
	log.Debugf("Dequeued task: %s, current size: %d", t.(*task.ScheduledTask).Task.Name(), bq.size)
	bq.fullWaitFn.Signal()
	return t.(*task.ScheduledTask)
}

func (bq *BoundedTaskQueue) Top() *task.ScheduledTask {
	log := logger.GetLogger()
	bq.lock.Lock()
	defer bq.lock.Unlock()

	for bq.size == 0 {
		log.Debug("Queue is empty, waiting for top task")
		bq.emptyWaitFn.Wait()
	}

	t := bq.pq.Top().(*task.ScheduledTask)
	log.Debugf("Top task: %s", t.Task.Name())
	return t
}
