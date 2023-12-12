// Algorithm for RRT* adapted from section 3.5.1 of:
// http://motion.cs.illinois.edu/RoboticSystems/MotionPlanningHigherDimensions.html

package rrtstar

import (
	"proj3-redesigned/robotpath"
)

// PathUpdateTask updates the path through the task of adding a milestone
type PathUpdate struct {
	path       *robotpath.Path
	mileStone  *robotpath.MileStone
	updateCost bool
	Done       chan any
}

// Create a new PathUpdateTask
func NewUpdate(path *robotpath.Path, updateCostInternally bool) *PathUpdate {
	return &PathUpdate{
		path:       path,
		mileStone:  nil,
		updateCost: updateCostInternally,
		Done:       make(chan any),
	}
}

// Run update according to the RRT* algorithm
func (task *PathUpdate) Run() {
	defer close(task.Done)
	task.mileStone = SamplePoint(task.path)
	Rewire(task.mileStone, task.path, task.updateCost)
}

// Get new milestone
func (task *PathUpdate) GetMileStone() *robotpath.MileStone {
	return task.mileStone
}
