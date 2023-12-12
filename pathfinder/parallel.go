package main

import (
	"proj3-redesigned/concurrent"
	"proj3-redesigned/robotpath"
	"proj3-redesigned/rrtstar"
)

// RunParallel runs the pathfinding algorithm in parallel
func RunParallel(configFile string, n int, threads int, strategy string) *robotpath.Path {
	// Read the configuration space from the input file
	path := robotpath.NewPath(configFile)

	// Initialize executor
	var executor concurrent.ExecutorService[rrtstar.PathUpdate, any]
	var updateCostInternally bool

	if strategy == "ws" {
		// Work stealing executor
		maxGrab := 100
		executor = concurrent.NewWorkStealingExecutor(threads, maxGrab)
		updateCostInternally = true

	} else if strategy == "bsp" {
		// BSP executor
		executor = concurrent.NewBSPExecutor(threads)
		updateCostInternally = false
	}

	// Populate the queues with tasks
	for i := 0; i < n; i++ {
		task := rrtstar.NewUpdate(path, updateCostInternally)
		executor.Submit(task)
	}

	// Execute
	executor.Execute()

	// Shutdown executor
	executor.Shutdown()

	return path
}
