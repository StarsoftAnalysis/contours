// mncontours.go

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

type PointT struct {
	x int
	y int
}

/* For now, points are just integer indexs into Pix
func (p PointT) Equal(p2 PointT) bool {
	return p.x == p2.x && p.y == p2.y
}
func (p PointT) Plus(p2 PointT) PointT {
	return PointT{p.x + p2.x, p.y + p2.y}
}
func (p PointT) Minus(p2 PointT) PointT {
	return PointT{p.x - p2.x, p.y - p2.y}
}

//var nilPoint = PointT{-1, -1}
*/

type ContourT []int // PointT

// Options and derived things
type OptsT struct {
	infile     string // from user
	width      int    // TODO ptr to img here instead?
	height     int    // pixels
	thresholds []int
	margin     float64 // from user
	paper      string  // from user
}

const white = 0xff
const black = 0x00

// Get the pixel value (0..255) at the given offset into the image
func getPix(imageData *image.NRGBA, i int) int {
	//return int(imageData.Pix[i*4]) // For now, just get the red part of the RGBA
	// Get grey:
	// Y = 0.299 R + 0.587 G + 0.114 B
	return int(math.Round(0.299*float64(imageData.Pix[i*4]) + 0.587*float64(imageData.Pix[i*4+1]) + 0.114*float64(imageData.Pix[i*4+2])))
}

var neighbourOffset = [8]PointT{
	{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0},
}

// Return a flag for each of 8 neighbours, set to true if the pixel is dark, i.e. within the shape.
// The neighbours are ordered like this (with y increasing downwards):
//
//	0  1  2
//	7     3
//	6  5  4
//
// Cells that are off the edge of the image are returned as false.
// Properly x,y version:
func neighboursWithinXY(imageData *image.NRGBA, width, height int, threshold int, p int) ([8]int, [8]bool) {
	var neighbours [8]int
	var within [8]bool = [8]bool{true, true, true, true, true, true, true, true}
	px := p % width
	py := p / width
	if px == 0 {
		// turn off left edge
		within[0] = false
		within[7] = false
		within[6] = false
	}
	if px == width-1 {
		// turn off right edge
		within[2] = false
		within[3] = false
		within[4] = false
	}
	if py == 0 {
		// turn off top edge
		within[0] = false
		within[1] = false
		within[2] = false
	}
	if py == height-1 {
		// turn off bottom edge
		within[6] = false
		within[5] = false
		within[4] = false
	}
	for i := 0; i < 8; i++ {
		neighbour := p + neighbourOffset[i].x + neighbourOffset[i].y*width
		// convert to x,y
		nx := neighbour % width
		ny := neighbour / width // integer division
		//fmt.Printf("nWXY new neighbour p=%v offset=%v width=%v neighbour=%v nx=%v ny=%v\n", p, neighbourOffset[i], width, neighbour, nx, ny)
		// check we're not off the edge of the image
		if nx >= 0 && nx < width && ny >= 0 && ny < height {
			if getPix(imageData, neighbour) >= threshold {
				within[i] = false
			}
		}
		neighbours[i] = neighbour // even if it's off the edge
	}
	//fmt.Printf("nWXY p=%v returning %v %v\n", p, neighbours, within)
	return neighbours, within
}

func traceContour(imageData *image.NRGBA, width, height int, threshold int, start int, svgF *SVGfile) ContourT {
	contour := make(ContourT, 1, 10)
	contour[0] = start
	var direction int = 3
	p := start
	//fmt.Printf("\ntC: width=%v start=%v\n", width, start)
	for true {
		neighbours, withins := neighboursWithinXY(imageData, width, height, threshold, p)
		// find the first neighbour starting from
		// the direction we came from
		var offset int = direction - 3 + 8
		/*
		   directions:
		     0  1  2
		     7     3
		     6  5  4
		   start indexes: -- is this what the -3 is for?
		     5  6  7
		     4     0
		     3  2  1
		*/
		nextP := -1
		for i := 0; i < 8; i++ {
			idx := (i + offset) % 8
			within := withins[idx]
			//fmt.Printf("tC loop: p=%v  offset=%v idx=%v ns=%v ws=%v\n", p, offset, idx, neighbours, withins)
			if within {
				direction = idx
				nextP = neighbours[idx]
				//fmt.Printf("tC: breaking with direction=%v nextP=%v\n", direction, nextP)
				break
			}
		}
		if nextP > width*height {
			fmt.Printf("tC: nextP=%v\n", nextP)
			panic("tC: p's out of range")
		}
		if nextP == -1 {
			// That's normal for a one-pixel shape
			if len(contour) > 1 {
				fmt.Printf("tC: !!!!!!!! finished loop without breaking\n")
				fmt.Printf("             p=%v nextP=%v contour=%v\n", p, nextP, contour)
			}
		}
		//fmt.Printf("tC: old p=%v  nextP=%v\n", p, nextP)
		p = nextP
		if p == start || p == -1 {
			break
		} else {
			contour = append(contour, p)
		}
	}
	if svgF != nil {
		// Single polygon -- assume the contour is closed
		svgF.polygon(contour, width)
	}
	return contour
}

func contourFinder(imageData *image.NRGBA, width, height int, threshold int, svgF *SVGfile) []ContourT {
	var contours = make([]ContourT, 0, 10)
	var imageLen = width * height
	seen := make([]bool, imageLen)
	var skipping = false
	for i := 0; i < imageLen; i++ {
		if getPix(imageData, i) < threshold {
			if !seen[i] && !skipping {
				contour := traceContour(imageData, width, height, threshold, i, svgF)
				contours = append(contours, contour)
				// this could be a _lot_ more efficient
				for _, c := range contour {
					seen[c] = true
				}
			}
			skipping = true
		} else {
			skipping = false
		}
	}
	return contours
}

/*
func writeSVG(contours []ContourT) {
	for _, contour := range contours {
	}
}
*/

func parseArgs(args []string) OptsT {
	var opts OptsT
	pf := pflag.NewFlagSet("contours", pflag.ExitOnError)
	pf.Float64VarP(&opts.margin, "margin", "m", 15, "Minimum margin (in mm).")
	pf.StringVarP(&opts.paper, "paper", "p", "A4L", "Paper size and orientation.  A4L | A4P | A3L | A3P.")
	pf.IntSliceVarP(&opts.thresholds, "threshold", "t", []int{128}, "Threshold levels, each 0..255")
	pf.SortFlags = false
	if args == nil {
		pf.Parse(os.Args[1:]) // don't pass program name
	} else {
		pf.Parse(args) // args passed as a string (for testing)
	}
	if pf.NArg() < 1 {
		fmt.Println("No input file name given")
		os.Exit(1)
	}
	opts.infile = pf.Arg(0)
	return opts
}

func main() {
	opts := parseArgs(nil)
	fmt.Printf("mncontours: processing '%s'\n", opts.infile)
	img, width, height, err := loadImage(opts.infile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	opts.width = width
	opts.height = height
	fmt.Printf("options: %#v\n", opts)
	optString := fmt.Sprintf("-mnc-t%sm%g%s", intsToString(opts.thresholds), opts.margin, opts.paper)
	ext := filepath.Ext(opts.infile)
	var svgF *SVGfile = new(SVGfile)
	svgF.openStart(strings.TrimSuffix(opts.infile, ext)+optString+".svg", opts)
	for t, threshold := range opts.thresholds {
		if svgF != nil {
			svgF.layer(t)
		}
		contours := contourFinder(img, opts.width, opts.height, threshold, svgF)
		fmt.Printf("%d contours found at threshold %d\n", len(contours), threshold)
	}
	svgF.stopSave()
}
