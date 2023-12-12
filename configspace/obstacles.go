package configspace

import (
	"image/color"
	"strconv"

	"github.com/fogleman/gg"
)

// Obstacle is an interface for obstacles in the configuration space
type Obstacle interface {
	SegmentCollision(*Point, *Point) bool
	Draw(*gg.Context)
}

// rectangleObstacle implements an Obstacle
type rectangleObstacle struct {
	pt *Point
	w  float32
	h  float32
}

// Creates a new rectangleObstacle Obstacle
func NewRectangleObstacle(config []string) Obstacle {
	x, _ := strconv.ParseFloat(config[0], 32)
	y, _ := strconv.ParseFloat(config[1], 32)
	h, _ := strconv.ParseFloat(config[2], 32)
	w, _ := strconv.ParseFloat(config[3], 32)

	return &rectangleObstacle{
		&Point{float32(x), float32(y)},
		float32(w),
		float32(h),
	}
}

// Detect obstacle's collision with the line segment described by the two points
func (r *rectangleObstacle) SegmentCollision(pt1 *Point, pt2 *Point) bool {
	// Get line segments that make up rectangle
	ll, lr := &Point{r.pt.X, r.pt.Y}, &Point{r.pt.X + r.w, r.pt.Y}
	ul, ur := &Point{r.pt.X, r.pt.Y + r.h}, &Point{r.pt.X + r.w, r.pt.Y + r.h}

	rectangleSegments := [][]*Point{
		{ll, ul}, {ul, ur},
		{ur, lr}, {lr, ll},
	}

	// Check each line segment
	for _, seg := range rectangleSegments {
		if Intersection(seg[0], seg[1], pt1, pt2) {
			return true
		}
	}
	return false
}

// Draw the obstacle onto the screen
func (r *rectangleObstacle) Draw(screen *gg.Context) {
	// Draw the image
	screen.SetColor(color.Black)
	screen.DrawRectangle(float64(r.pt.X), float64(r.pt.Y), float64(r.w), float64(r.h))
	screen.Fill()
}
