package configspace

// Unit testing for Obstacles.go. Tests the following functions:
// NewPoint
// NewRectangleObstacle
// SegmentCollision
//

import (
	"testing"
)

// Test NewPoint
func TestNewPoint(t *testing.T) {
	pt := Point{1.0, 2.0}
	if pt.X != 1.0 || pt.Y != 2.0 {
		t.Error("NewPoint failed")
	}
}

// Test NewRectangleObstacle
func TestNewRectangleObstacle(t *testing.T) {
	rect := NewRectangleObstacle([]string{"0.0", "0.0", "0.5", "0.5"})
	if rect.pt.X != 0.0 || rect.pt.Y != 0.0 || rect.w != 0.5 || rect.h != 0.5 {
		t.Error("NewRectangleObstacle failed")
	}
}

// Test SegmentCollision
func TestSegmentCollisionDiagonalThrough(t *testing.T) {
	rect := NewRectangleObstacle([]string{"1.0", "1.0", "0.5", "0.5"})
	pt1 := &Point{0.0, 0.0}
	pt2 := &Point{2.0, 2.0}
	if !rect.SegmentCollision(pt1, pt2) {
		t.Error("SegmentCollision failed")
	}
}

func TestSegmentCollisionCollinear(t *testing.T) {
	rect := NewRectangleObstacle([]string{"0.0", "0.5", "0.5", "0.5"})
	pt1 := &Point{1.0, 0.5}
	pt2 := &Point{1.5, 0.5}
	if rect.SegmentCollision(pt1, pt2) {
		t.Error("SegmentCollision failed")
	}
}

func TestSegmentCollisionEndPoint(t *testing.T) {
	rect := NewRectangleObstacle([]string{"0.0", "0.5", "1.0", "1.0"})
	pt1 := &Point{0.5, 0.0}
	pt2 := &Point{0.5, 0.5}
	if !rect.SegmentCollision(pt1, pt2) {
		t.Error("SegmentCollision failed")
	}
	pt1 = &Point{0.5, 0.0}
	pt2 = &Point{0.5, 0.25}
	if rect.SegmentCollision(pt1, pt2) {
		t.Error("SegmentCollision failed")
	}
}
