package rrtstar

import (
	"proj3-redesigned/robotpath"
)

// Rewiring of the RRT* algorithm, assumes milestone that is passed is randomly
// drawn and valid w.r.t. obstacles in the configuration space, see SamplePoint()
func Rewire(ms *robotpath.MileStone, path *robotpath.Path, doCostUpdate bool) {
	// Rewire the tree to account for the new milestone
	rewirePath(ms, path)

	// Check if milestone is most optimal path to goal
	checkSuccess := false
	for !checkSuccess {
		distBetweenGoal := robotpath.Distance(ms.Point, path.Goal.Point)
		goalCost := path.Goal.Cost

		if distBetweenGoal < path.Config.Visibility &&
			(goalCost == 0.0 || ms.Cost+distBetweenGoal < goalCost) {
			checkSuccess = tryRewire(path.Goal, ms, goalCost, distBetweenGoal, path)

		} else {
			checkSuccess = true
		}
	}

	// Run cost update
	if doCostUpdate {
		ms.UpdateChildrenCost()
	}
}

// Rewire the tree to account for the new MileStone
func rewirePath(ms *robotpath.MileStone, path *robotpath.Path) {
	// Find 10 nearest neighbors in the path
	nHood := path.GetNN(ms, 10)

	// Check if each neighbor requires re-wireing
	for _, n := range nHood {
		checkSuccess := false

		for !checkSuccess {
			// Determine relative distances between points
			distBetween := robotpath.Distance(ms.Point, n.Point)
			msCost, neighborCost := ms.Cost, n.Cost
			distThroughNew := msCost + distBetween
			distToNew := neighborCost + distBetween

			if distThroughNew < neighborCost {
				// Shorter path from new milestone to neighbor
				checkSuccess = tryRewire(n, ms, neighborCost, distBetween, path)

			} else if distToNew < msCost {
				// Shorter path from neighbor to new milestone
				checkSuccess = tryRewire(ms, n, msCost, distBetween, path)

			} else {
				// No re-wireing needed
				checkSuccess = true
			}
		}
	}
}

// Attempts to rewire two points, returns false if re-attempt necessary due to
// changes to newChild's parent, sequential program will always return true
func tryRewire(newChild *robotpath.MileStone, newParent *robotpath.MileStone,
	childCost float32, dist float32, path *robotpath.Path,
) bool {
	// Check if points visible to another
	if path.Config.Visible(newChild.Point, newParent.Point) {
		// Attempt to set new parent
		if !newChild.SetParent(newParent, childCost, dist) {
			return false
		}
	}
	return true
}
