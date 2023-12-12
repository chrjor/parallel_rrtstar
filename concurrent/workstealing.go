package concurrent

import (
	"math/rand"
	"proj3-redesigned/deque"
	"proj3-redesigned/rrtstar"
	"sync"
)

// Work stealing executor service
type WorkStealingExecutor struct {
	workers   []*Worker        // The workers in the pool
	wg        sync.WaitGroup   // WaitGroup for workers
	threshold int              // Threshold for grabbing from the queue
	shutdown  chan interface{} // Channel for shutdown
	tasks     int              // Counter for initializing the queues
}

// Worker struct
type Worker struct {
	queue deque.DEQue[rrtstar.PathUpdate]
}

// NewWorkStealingExecutor returns an ExecutorService that is implemented using the
// work-stealing algorithm. Capacity is the number of goroutines in the pool and
// threshold is the number of items that a goroutine in the pool can grab from the
// executor in one time period
func NewWorkStealingExecutor(capacity, threshold int,
) ExecutorService[rrtstar.PathUpdate, any] {
	// Create worker array
	var workers []*Worker
	for i := 0; i < capacity; i++ {
		var worker Worker
		worker.queue = deque.NewUnboundedDEQue[rrtstar.PathUpdate]()
		workers = append(workers, &worker)
	}
	// Create executor
	executor := &WorkStealingExecutor{
		workers:   workers,
		wg:        sync.WaitGroup{},
		threshold: threshold,
		shutdown:  make(chan interface{}),
	}

	return executor
}

// Submits a task to the executor
func (e *WorkStealingExecutor) Submit(task *rrtstar.PathUpdate) Future[any] {
	e.workers[e.tasks%len(e.workers)].queue.PushBottom(task)
	e.wg.Add(1)
	e.tasks++
	return &RunnableFuture{Done: task.Done}
}

func (e *WorkStealingExecutor) Execute() {
	// Run the workers
	for worker := 0; worker < len(e.workers); worker++ {
		go e.runWorker(worker)
	}
}

// runWorkStealer is the main worker instructions for the worker stealing routine
func (e *WorkStealingExecutor) runWorker(me int) {
	for {
		select {
		case <-e.shutdown:
			return
		default:
			// If local queue is empty, steal from a random worker
			if e.workers[me].queue.IsEmpty() {
				randSteal := rand.Intn(len(e.workers))
				for randSteal == me {
					randSteal = rand.Intn(len(e.workers))
				}
				for i := 0; i < e.threshold; i++ {
					task := e.workers[randSteal].queue.PopTop()
					if task != nil {
						e.workers[me].queue.PushBottom(task)
					} else {
						break
					}
				}
			}

			task := e.workers[me].queue.PopBottom()
			if task != nil {
				task.Run()
				e.wg.Done()
			}
		}
	}
}

// Shuts down the executor
func (e *WorkStealingExecutor) Shutdown() {
	e.wg.Wait()
	close(e.shutdown)
}
