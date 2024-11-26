package container

import (
	"github.com/Adarsh-kumar/go-concurrency/pkg/task"
)

type PriorityQueue []*task.ScheduledTask

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].GetScheduledTime().Before(pq[j].GetScheduledTime())
}

func (pq *PriorityQueue) Push(t interface{}) {
	*pq = append(*pq, t.(*task.ScheduledTask))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (pq *PriorityQueue) Top() interface{} {
	return (*pq)[0]
}
