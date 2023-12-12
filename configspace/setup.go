package configspace

import (
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

// ConfigSpace is a struct representing the configuration space
type Config struct {
	Start      *Point     // Start point
	Goal       *Point     // Goal point
	Visibility float32    // Visibility radius
	Obstacles  []Obstacle // Obstacles in the configuration space
	WinHeight  float32    // Window height
	WinWidth   float32    // Window width
}

// Point is a general struct used for points
type Point struct {
	X float32
	Y float32
}

// Create a new configuration space from a config file
func NewConfigSpace(configPath string) *Config {

	// Initialize space's variables
	var winWidth, winHeight, radius float64
	var start, goal *Point
	var obstacles []Obstacle

	// Parse config file
	config := ReadFile(configPath)

	for _, line := range config {

		line := strings.Split(line, ",")

		if line[0] == "window" {
			winHeight, _ = strconv.ParseFloat(line[1], 32)
			winWidth, _ = strconv.ParseFloat(line[2], 32)

		} else if line[0] == "visibility" {
			radius, _ = strconv.ParseFloat(line[1], 32)

		} else if line[0] == "start" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			start = &Point{float32(x), float32(y)}

		} else if line[0] == "goal" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			goal = &Point{float32(x), float32(y)}

		} else if line[0] == "rectangle" {
			obstacles = append(obstacles, NewRectangleObstacle(line[1:]))
		}
	}

	return &Config{
		Start:      start,
		Goal:       goal,
		Visibility: float32(radius),
		Obstacles:  obstacles,
		WinHeight:  float32(winHeight),
		WinWidth:   float32(winWidth),
	}
}

// NewPoint creates a new Point
func (c *Config) NewPoint(x, y float32) *Point {
	return &Point{x, y}
}

// Check if a new path branch (line segment) is not obstructed by any obstacle
func (c *Config) Visible(pt1 *Point, pt2 *Point) bool {
	for _, o := range c.Obstacles {
		if o.SegmentCollision(pt1, pt2) {
			return false
		}
	}
	return true
}

// Draw the configuration space
func (c *Config) Draw(screen *gg.Context) {
	for _, o := range c.Obstacles {
		o.Draw(screen)
	}
}
