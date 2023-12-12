package robotpath

import (
	"container/heap"
	"image/color"
	"math"
	"proj3-redesigned/configspace"
	"sync"

	"github.com/fogleman/gg"
)

// Path is the struct that oversees the path planning process
type Path struct {
	Config     *configspace.Config // Configuration space
	Goal       *MileStone          // Goal milestone
	Start      *MileStone          // Start milestone
	milestones []*MileStone        // Milestone array of nodes in the tree
	rw         sync.RWMutex        // Mutex to update milestone array
}

// Create a new Path object and set its configuration space
func NewPath(configPath string) *Path {
	path := Path{
		Config:     configspace.NewConfigSpace(configPath),
		milestones: make([]*MileStone, 0),
		rw:         sync.RWMutex{},
	}

	path.Start = NewMileStone(path.Config.Start)
	path.Goal = NewMileStone(path.Config.Goal)

	// We begin with a single start milestone
	path.AddPoint(path.Start)

	return &path
}

// Get k-nearest neighbors of a milestone
func (path *Path) GetNN(ms *MileStone, k int) []*MileStone {

	var neighborhood NeighborHeap

	// Lock the path's milestone array from writers and process nearest neighbor heap
	path.rw.RLock()
	for _, oldMs := range path.milestones {
		dist := Distance(oldMs.Point, ms.Point)
		neighbor := NewNeighborItem(oldMs, dist)
		heap.Push(&neighborhood, neighbor)
	}
	path.rw.RUnlock()

	// Extract k-nearest neighbors from heap
	var neighbors []*MileStone
	numNeighbors := math.Min(float64(k), float64(len(neighborhood)))
	for i := 0; i < int(numNeighbors); i++ {
		n := heap.Pop(&neighborhood).(*NeighborItem).Neighbor
		if n != ms {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

// Add milestone to path
func (path *Path) AddPoint(newMs *MileStone) {
	// Write-lock path's milestone array before appending
	path.rw.Lock()
	defer path.rw.Unlock()
	path.milestones = append(path.milestones, newMs)
}

// Get minimum distance to goal
func (path *Path) DistToGoal() float32 {
	return path.Goal.Cost
}

// Draw the path and configuration space
func (path *Path) Draw(screen *gg.Context) {

	lightBlue := color.RGBA{R: 173, G: 216, B: 230, A: 255}
	darkRed := color.RGBA{R: 139, G: 0, B: 0, A: 255}
	darkGreen := color.RGBA{R: 0, G: 100, B: 0, A: 255}

	// Draw obstacles
	path.Config.Draw(screen)

	// Draw path tree
	var treeDraw func(*MileStone)
	treeDraw = func(lastPt *MileStone) {
		lastPt.Children.Range(func(key, value any) bool {
			child := value.(*MileStone)
			screen.SetLineWidth(5.0)
			screen.SetColor(lightBlue)
			screen.DrawLine(float64(lastPt.Point.X), float64(lastPt.Point.Y),
				float64(child.Point.X), float64(child.Point.Y))
			screen.Stroke()
			treeDraw(child)
			return true
		})
	}
	treeDraw(path.Start)

	// Draw optimal path and Start point
	screen.SetColor(darkGreen)
	screen.DrawPoint(float64(path.Start.Point.X), float64(path.Start.Point.Y), 5.0)
	screen.Fill()
	ms := path.Goal
	for ms.Parent != nil {
		screen.SetLineWidth(6.0)
		screen.SetColor(darkGreen)
		screen.DrawLine(float64(ms.Point.X), float64(ms.Point.Y),
			float64(ms.Parent.Point.X), float64(ms.Parent.Point.Y))
		screen.Stroke()
		ms = ms.Parent
	}

	// Draw Start and Goal points
	screen.SetColor(darkRed)
	screen.DrawPoint(float64(path.Goal.Point.X), float64(path.Goal.Point.Y), 5.0)
	screen.Fill()
}

// Calculate the distance between two points in the configuration space
func Distance(pt1 *configspace.Point, pt2 *configspace.Point) float32 {
	base_sq := math.Pow(float64(pt1.X-pt2.X), 2)
	height_sq := math.Pow(float64(pt1.Y-pt2.Y), 2)

	return float32(math.Sqrt(base_sq + height_sq))
}
