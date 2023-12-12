// Algorithm used for line-segment intersection detection adapted from CLRS
// Chapter 33 and a corresponding 2D approach at:
// https://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/

package configspace

import (
	"math"
)

// Check if two 2D line segments intersect using their end points
func Intersection(p1 *Point, q1 *Point, p2 *Point, q2 *Point) bool {
	o1 := SegmentOrientation(p1, q1, p2)
	o2 := SegmentOrientation(p1, q1, q2)
	o3 := SegmentOrientation(p2, q2, p1)
	o4 := SegmentOrientation(p2, q2, q1)

	// Check the general case, else check when collinearity present
	if (o1 != o2) && (o3 != o4) {
		return true
	} else if (o1 == 0 && OnSegment(p1, p2, q1)) ||
		(o2 == 0 && OnSegment(p1, q2, q1)) ||
		(o3 == 0 && OnSegment(p2, p1, q2)) ||
		(o4 == 0 && OnSegment(p2, q1, q2)) {
		return true
	}
	return false
}

// Calculate the orientation of three 2D points
func SegmentOrientation(pt1 *Point, pt2 *Point, pt3 *Point) int {
	orient := ((pt2.Y - pt1.Y) * (pt3.X - pt2.X)) -
		((pt2.X - pt1.X) * (pt3.Y - pt2.Y))

	if orient > 0 {
		// Clockwise
		return 1

	} else if orient < 0 {
		// Counterclockwise
		return 2
	}
	// Collinear
	return 0
}

// If points are collinear, detect if pt2 lies on the line pt1pt3
func OnSegment(pt1 *Point, pt2 *Point, pt3 *Point) bool {
	x1, y1 := float64(pt1.X), float64(pt1.Y)
	x2, y2 := float64(pt2.X), float64(pt2.Y)
	x3, y3 := float64(pt3.X), float64(pt3.Y)
	if (x2 <= math.Max(x1, x3)) &&
		(x2 >= math.Min(x1, x3)) &&
		(y2 <= math.Max(y1, y3)) &&
		(y2 >= math.Min(y1, y3)) {
		return true
	}
	return false
}
