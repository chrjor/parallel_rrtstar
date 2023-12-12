package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"proj3-redesigned/robotpath"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)


// Usage statement
const usage = "\nUsage:	go run proj3-redesigned/pathfinder <bench|sim> <samples> <input_file> [ws|bsp] [threads] \n\n" +
	"Mandatory Arguments:\n" +
	"- <bench|sim>:		benchmark mode or simulation mode which outputs an image\n" +
	"- <samples>:		number of samples drawn to find the path\n" +
	"- <input_file>:		file for configuration space setup\n\n" +
	"Optional Arguments:\n" +
	"- [ws|bsp]:		work stealing or bulk synchronous parallel scheduling\n" +
	"- [threads]:		number of threads when selecting parallized version\n" +
	"\nNote: Omit [ws|bsp] and [threads] for sequential program\n\n" +
	"Examples:\n" +
	"- Sequental:	go run proj3-redesigned/pathfinder bench 1000 data/maze.txt\n" +
	"- Parallel:	go run proj3-redesigned/pathfinder bench 1000 data/maze.txt ws 4\n"

func main() {
	// Check for correct number of command line arguments
	if !(len(os.Args) == 4 || len(os.Args) == 6) {
		fmt.Println(usage)
		return
	}
	// Parse command line arguments
	mode := os.Args[1]
	if mode != "bench" && mode != "sim" {
		fmt.Println(usage)
		return
	}
	sampleSize, _ := strconv.Atoi(os.Args[2])
	inputPath := os.Args[3]
	var strategy string
	threads := 1
	if len(os.Args) == 6 {
		strategy = os.Args[4]
		if strategy != "ws" && strategy != "bsp" {
			fmt.Println(usage)
			return
		}
		threads, _ = strconv.Atoi(os.Args[5])
	}

	// Start benchmark timer
	start := time.Now()

	// Run program
	var output *robotpath.Path
	if threads == 1 {
		// Sequential program
		output = RunSequential(inputPath, sampleSize)
	} else {
		// Parallel program
		output = RunParallel(inputPath, sampleSize, threads, strategy)
	}

	// Print benchmark time
	end := time.Since(start).Seconds()
	fmt.Printf("%.2f\n", end)

	if mode == "sim" {
		// Create image of simulation results
		img := image.NewRGBA(image.Rect(0, 0, int(output.Config.WinWidth), int(output.Config.WinHeight)))
		screen := gg.NewContextForRGBA(img)
		screen.SetColor(color.White)
		screen.Clear()
		output.Draw(screen)

		// Write image to file
		pathName := fmt.Sprintf("data/output/maze_%d.jpg", sampleSize)
		f, _ := os.Create(pathName)
		defer f.Close()
		jpeg.Encode(f, img, &jpeg.Options{Quality: 100})

		// Print output
		fmt.Println("Goal distance: ", output.DistToGoal())
		fmt.Println("Image created.")
	}
}
