package rrtstar

import (
	"math"
	"math/rand"
	"proj3-redesigned/robotpath"
)

// SamplePoint samples a valid random point in the configuration space
func SamplePoint(path *robotpath.Path) *robotpath.MileStone {
	var ms *robotpath.MileStone
	ms = nil

	// Sample until valid milestone created
	for ms == nil {
		randX := rand.Float32() * float32(path.Config.WinWidth)
		randY := rand.Float32() * float32(path.Config.WinHeight)
		ms = tryPathExtend(robotpath.NewMileStone(path.Config.NewPoint(randX, randY)), path)
	}
	return ms
}

// Extend the path from randomly drawn point to the nearest point in the tree
func tryPathExtend(ms *robotpath.MileStone, path *robotpath.Path) *robotpath.MileStone {
	// Find nearest neighbor to the sampled point
	nHood := path.GetNN(ms, 1)
	nearest := nHood[0]

	// Extend path to nearest neighbor and check if new position valid, if not
	// restart the process by returning nil
	extend(ms, nearest, path.Config.Visibility)
	if !path.Config.Visible(ms.Point, nearest.Point) {
		return nil
	}

	// Add the point to the path plan
	newDist := robotpath.Distance(nearest.Point, ms.Point)
	ms.SetParent(nearest, 0.0, newDist)
	path.AddPoint(ms)

	return ms
}

// Set milestone's new location as distance from its nearest neighbor to the
// closest point in the direction of its current position
func extend(ms *robotpath.MileStone, nearest *robotpath.MileStone, radius float32) {
	length := robotpath.Distance(ms.Point, nearest.Point)
	vis := float32(math.Min(float64(radius), float64(length)))

	ms.Point.X = nearest.Point.X + (ms.Point.X-nearest.Point.X)*vis/length
	ms.Point.Y = nearest.Point.Y + (ms.Point.Y-nearest.Point.Y)*vis/length
}
