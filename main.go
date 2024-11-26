package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Adarsh-kumar/go-concurrency/pkg/executor"
)

// Example implementation of a Task
type SimpleTask struct {
	name string
	wg   *sync.WaitGroup
}

func (t *SimpleTask) Run() error {
	defer t.wg.Done() // Signal task completion to the WaitGroup
	fmt.Printf("[%s] Task %s is running\n", time.Now().Format("15:04:05"), t.name)
	time.Sleep(500 * time.Millisecond) // Simulate some work
	return nil
}

func (t *SimpleTask) Name() string {
	return t.name
}

// Function to test the scheduler
func testScheduler() {
	var wg sync.WaitGroup

	// Use the number of CPUs as the maximum number of threads
	numProcessors := runtime.NumCPU()
	scheduledExecutor := executor.NewScheduledExecutor(numProcessors)
	fmt.Printf("Scheduler initialized with %d processors\n", numProcessors)

	// Schedule a task with a fixed delay
	wg.Add(1) // Increment the WaitGroup counter
	fixedDelayTask := &SimpleTask{name: "Fixed-Delay Task", wg: &wg}
	scheduledExecutor.ScheduleTaskWithFixedDelay(fixedDelayTask, 10*time.Second)

	// Simulate submission of a new task after a delay
	wg.Add(1)
	go func() {
		time.Sleep(3 * time.Second)
		newTask := &SimpleTask{name: "New Task During Wait", wg: &wg}
		fmt.Printf("[%s] Submitting a new task after 3 seconds\n", time.Now().Format("15:04:05"))
		scheduledExecutor.ScheduleTask(newTask)
	}()

	// Simulate submission of another task with a shorter delay
	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		shortTask := &SimpleTask{name: "Short Task", wg: &wg}
		fmt.Printf("[%s] Submitting a short task after 5 seconds\n", time.Now().Format("15:04:05"))
		scheduledExecutor.ScheduleTask(shortTask)
	}()

	// Wait for all tasks to complete
	wg.Wait()
	fmt.Println("Test completed. Stopping executor...")
}

func main() {
	// Call the test function
	testScheduler()
}
