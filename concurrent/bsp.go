package concurrent

import (
	"math"
	"proj3-redesigned/robotpath"
	"proj3-redesigned/rrtstar"
	"sync"
	"sync/atomic"
)

// BSP executor service
type BSPExecutor struct {
	ctx *bspContext 	// BSP context
	wg  sync.WaitGroup  // WaitGroup for tracking task completion
}

// BSP context
type bspContext struct {
	numWorkers   int32 					// Number of workers
	numSyncing   atomic.Int32 			// Number of workers syncing
	taskBuffer   []*rrtstar.PathUpdate 	// All remaining tasks
	curWork      []*rrtstar.PathUpdate 	// Worker's current task
	syncMessages []*robotpath.MileStone // Messages used for synchronization
	cond         sync.Cond 				// Condition variable for synchronization
	shutdown     chan interface{} 		// Channel for shutdown
}

// NewBSPExecutor returns an ExecutorService that is implemented using the BSP
// scheduling strategy
func NewBSPExecutor(threads int) ExecutorService[rrtstar.PathUpdate, any] {
	// Create BSP context
	context := bspContext{
		numWorkers:   int32(threads),
		numSyncing:   atomic.Int32{},
		taskBuffer:   make([]*rrtstar.PathUpdate, 0),
		curWork:      make([]*rrtstar.PathUpdate, threads),
		syncMessages: make([]*robotpath.MileStone, threads),
		cond:         *sync.NewCond(&sync.Mutex{}),
		shutdown:     make(chan interface{}),
	}
	// Create executor
	executor := &BSPExecutor{
		ctx: &context,
		wg:  sync.WaitGroup{},
	}
	return executor
}

// Submits a task to the executor
func (e *BSPExecutor) Submit(task *rrtstar.PathUpdate) Future[any] {
	e.ctx.taskBuffer = append(e.ctx.taskBuffer, task)
	return &RunnableFuture{Done: task.Done}
}

// Executes the executor
func (e *BSPExecutor) Execute() {
	// Define worker loop
	runBSPWorker := func(id int, ctx *bspContext) {
		defer e.wg.Done()
		for {
			select {
			case <-e.ctx.shutdown:
				return
			default:
				// Sync with other threads
				ctx.Sync(id)
				// Execute work
				if ctx.curWork[id] != nil {
					ctx.curWork[id].Run()
					ctx.syncMessages[id] = ctx.curWork[id].GetMileStone()
				}
			}
		}
	}
	// Run the workers
	for worker := 0; worker < int(e.ctx.numWorkers); worker++ {
		e.wg.Add(1)
		go runBSPWorker(worker, e.ctx)
	}
}

// Shuts down the executor
func (e *BSPExecutor) Shutdown() {
	e.wg.Wait()
}

// Synchronizes the workers in between the BSP steps
func (ctx *bspContext) Sync(id int) {
	// Check if thread last to finish task, if so run update() to synchronize
	ctx.cond.L.Lock()
	if ctx.numSyncing.Load() < ctx.numWorkers-1 {
		ctx.numSyncing.Add(1)
		ctx.cond.Wait()
		ctx.numSyncing.Add(-1)
	} else {
		ctx.update()
		ctx.cond.Broadcast()
	}
	ctx.cond.L.Unlock()
}

// Updates the BSP context in between steps
func (ctx *bspContext) update() {
	// Update the cost of the new milestones
	for _, newMilestone := range ctx.syncMessages {
		if newMilestone != nil {
			newMilestone.UpdateChildrenCost()
		}
	}

	// Update current work for each worker	
	numTasks := len(ctx.taskBuffer)
	for i := 0; i < int(ctx.numWorkers); i++ {
		if numTasks > i {
			ctx.curWork[i] = ctx.taskBuffer[numTasks-i-1]
		} else {
			ctx.curWork[i] = nil
		}
	}

	// Update task buffer
	newEnd := numTasks - int(math.Min(float64(numTasks), float64(ctx.numWorkers)))

	ctx.taskBuffer = ctx.taskBuffer[:newEnd]

	if len(ctx.taskBuffer) == 0 {
		close(ctx.shutdown)
	}
}
