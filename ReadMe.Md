# Concurrent Scheduled Executor

This project implements a **Concurrent Scheduled Executor** in Go that supports scheduling tasks with different policies. It is backed by a **bounded blocking priority queue**, ensuring efficient and thread-safe task execution with a fixed number of worker threads.

## Features

- **Task Scheduling Policies**:
  - **One-shot Task**: Executes a task at a specific time.
  - **Fixed Delay Task**: Executes a task repeatedly with a fixed delay between each execution.
  - **Fixed Rate Task**: Executes a task repeatedly at fixed intervals, regardless of task execution time.

- **Bounded Blocking Priority Queue**:
  - Maintains tasks in order of their scheduled execution time.
  - Blocks enqueue operations when the queue is full, ensuring bounded resource usage.
  - Blocks dequeue operations when the queue is empty, waiting for tasks to become available.

- **Task Interface**:
  - Defines a generic task structure with methods for execution and identification.

---

## Components 

1. **Task Interface**:
   Any task to be scheduled must implement the `Task` interface:

   ```go
   type Task interface {
       Run() error    // Defines the task's execution logic
       Name() string  // Returns the task's name for logging/debugging
   }
2. **Executor**:

    The `ScheduledExecutor` provides methods to schedule tasks with different policies:

    - **`ScheduleTask`**: For one-shot tasks.
    - **`ScheduleTaskWithFixedDelay`**: For tasks with a fixed delay.
    - **`ScheduleTaskAtFixedRate`**: For tasks at a fixed rate.

---

3. **Bounded Blocking Priority Queue**:

    Used to store the tasks, ordered by their next scheduled execution time. The priority queue is essentially a concurrent bounded blocking queue, i.e The producer (task inflow) and the 
    consumer (scheduler) will wait without holding the lock and will be notified when the condition is met. 
  
## Example Usage

### Example with a sample task.

 ```go
import (
    "github.com/Adarsh-kumar/go-concurrency/pkg/executor"
    "sync"
    "time"
)

// Example Task Implementation
type ExampleTask struct {
    name string
    wg   *sync.WaitGroup
}

func (t *ExampleTask) Run() error {
    defer t.wg.Done()
    fmt.Printf("[%s] Task %s is running\n", time.Now().Format("15:04:05"), t.name)
    time.Sleep(500 * time.Millisecond) // Simulate work
    return nil
}

func (t *ExampleTask) Name() string {
    return t.name
}

// Start the scheduler and submit the task 
func main() {
    var wg sync.WaitGroup

    // Create a new executor with threads equal to the number of processors
    numProcessors := runtime.NumCPU()
    scheduler := executor.NewScheduledExecutor(numProcessors)

    // Schedule a one-shot task
    wg.Add(1)
    oneShotTask := &ExampleTask{name: "One-Shot Task", wg: &wg}
    scheduler.ScheduleTask(oneShotTask)

    // Schedule a fixed delay task
    wg.Add(1)
    fixedDelayTask := &ExampleTask{name: "Fixed-Delay Task", wg: &wg}
    scheduler.ScheduleTaskWithFixedDelay(fixedDelayTask, 5*time.Second)

    // Schedule a fixed rate task
    wg.Add(1)
    fixedRateTask := &ExampleTask{name: "Fixed-Rate Task", wg: &wg}
    scheduler.ScheduleTaskAtFixedRate(fixedRateTask, 2*time.Second)

    // Wait for tasks to complete
    wg.Wait()
    fmt.Println("All tasks completed.")
}
```

### Output logs
```bash
{"level":"info","ts":1732636477.4344504,"caller":"executor/scheduled_executor.go:23","msg":"Initializing ScheduledExecutor"}
Scheduler initialized with 12 processors
{"level":"info","ts":1732636477.4346466,"caller":"executor/scheduled_executor.go:49","msg":"Scheduling task with fixed delay: Fixed-Delay Task, delay: 10s"}
{"level":"info","ts":1732636477.4347808,"caller":"executor/scheduled_executor.go:93","msg":"Next task (Fixed-Delay Task) scheduled in 9.999916591s"}
{"level":"info","ts":1732636477.4348853,"caller":"executor/scheduled_executor.go:93","msg":"Next task (Fixed-Delay Task) scheduled in 9.999790563s"}
[15:54:40] Submitting a new task after 3 seconds
{"level":"info","ts":1732636480.4351103,"caller":"executor/scheduled_executor.go:36","msg":"Scheduling one-shot task: New Task During Wait"}
{"level":"info","ts":1732636480.4352272,"caller":"executor/scheduled_executor.go:88","msg":"Executing task: New Task During Wait"}
{"level":"info","ts":1732636480.4352493,"caller":"executor/scheduled_executor.go:93","msg":"Next task (Fixed-Delay Task) scheduled in 6.999428946s"}
{"level":"info","ts":1732636480.4352984,"caller":"executor/scheduled_executor.go:122","msg":"Running task: New Task During Wait"}
[15:54:40] Task New Task During Wait is running
{"level":"info","ts":1732636480.9359992,"caller":"executor/scheduled_executor.go:127","msg":"Task New Task During Wait completed successfully"}
[15:54:42] Submitting a short task after 5 seconds
{"level":"info","ts":1732636482.4356158,"caller":"executor/scheduled_executor.go:36","msg":"Scheduling one-shot task: Short Task"}
{"level":"info","ts":1732636482.4356651,"caller":"executor/scheduled_executor.go:88","msg":"Executing task: Short Task"}
{"level":"info","ts":1732636482.4357069,"caller":"executor/scheduled_executor.go:93","msg":"Next task (Fixed-Delay Task) scheduled in 4.998982194s"}
{"level":"info","ts":1732636482.4357557,"caller":"executor/scheduled_executor.go:122","msg":"Running task: Short Task"}
[15:54:42] Task Short Task is running
{"level":"info","ts":1732636482.9365015,"caller":"executor/scheduled_executor.go:127","msg":"Task Short Task completed successfully"}
{"level":"info","ts":1732636487.4390998,"caller":"executor/scheduled_executor.go:97","msg":"Scheduled time reached for task: Fixed-Delay Task"}
{"level":"info","ts":1732636487.4392138,"caller":"executor/scheduled_executor.go:88","msg":"Executing task: Fixed-Delay Task"}
{"level":"info","ts":1732636487.439266,"caller":"executor/scheduled_executor.go:122","msg":"Running task: Fixed-Delay Task"}
[15:54:47] Task Fixed-Delay Task is running
{"level":"info","ts":1732636487.9400115,"caller":"executor/scheduled_executor.go:127","msg":"Task Fixed-Delay Task completed successfully"}
Test completed. Stopping executor...




