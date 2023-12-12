package main

import (
	"proj3-redesigned/robotpath"
	"proj3-redesigned/rrtstar"
)

// RunSequential runs the pathfinding algorithm sequentially
func RunSequential(configFile string, n int) *robotpath.Path {

	// Read the configuration space from the input file and create new path
	path := robotpath.NewPath(configFile)

	// Make n updates to the path using the RRT* algorithm
	for i := 0; i < n; i++ {
		task := rrtstar.NewUpdate(path, true)
		task.Run()
	}

	return path
}
