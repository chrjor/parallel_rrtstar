package robotpath

import (
	"proj3-redesigned/configspace"
	"sync"
)

// MileStone represents a node in the path plan
type MileStone struct {
	Point    *configspace.Point // Point coordinates
	Parent   *MileStone         // Parent milestone
	dist     float32            // Distance to parent
	Children sync.Map           // Child milestoned
	costLock sync.Mutex         // Lock for updating cost
	Cost     float32            // Cost (distance to start Milestone)
}

// Create a new MileStone using a point
func NewMileStone(pt *configspace.Point) *MileStone {
	return &MileStone{
		Point:    pt,
		Parent:   nil,
		Children: sync.Map{},
		costLock: sync.Mutex{},
		Cost:     0.0,
	}
}

// Set the parent of a milestone. The milestone locks it's cost in order to
// update its parent and cost
func (ms *MileStone) SetParent(newParent *MileStone, curCost float32, dist float32) bool {

	// Block any cost updates while switching parents
	ms.costLock.Lock()
	defer ms.costLock.Unlock()

	// A change has been made asynchronously which invalidates this attempt
	if curCost != ms.Cost {
		return false
	}

	// Parent link updated before cost so any cost updates to the parent wait on lock
	if ms.Parent != nil {
		ms.Parent.removeChild(ms)
	}
	ms.Parent = newParent
	ms.Parent.setChild(ms)

	// Get new parent's cost and update with dist between
	ms.Cost = ms.Parent.Cost + dist
	ms.dist = dist

	return true
}

// Add a child to a milestone
func (ms *MileStone) setChild(child *MileStone) {
	ms.Children.Store(child, child)
}

// Remove a child from a milestone
func (ms *MileStone) removeChild(child *MileStone) {
	ms.Children.Delete(child)
}

// Sets the cost of a milestone by checking the parent's cost and the dist
// between them, this will change if the parent's cost was updated
func (ms *MileStone) setCost() {

	ms.costLock.Lock()
	defer ms.costLock.Unlock()

	ms.Cost = ms.Parent.Cost + ms.dist
}

// Update the cost of all descendents of a milestone
func (ms *MileStone) UpdateChildrenCost() {
	ms.Children.Range(func(key, value any) bool {
		child := value.(*MileStone)
		child.setCost()
		child.UpdateChildrenCost()
		return true
	})
}
