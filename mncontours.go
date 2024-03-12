// mncontours.go

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

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
	infile string // from user
	width  int    // TODO ptr to img here instead?
	height int    // pixels
}

var opts OptsT

const white = 0xff
const black = 0x00
const threshold = 0x7f

// Relative neighbours
//
//	0   1   2   y down
//	7       3
//	6   5   4
//
// relNeighbours := make([]PointT, 8);
/*
var relNeighbours = []PointT{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{+1, -1}, {1, 0}, {1, 1},
}
*/

// list of neighbours to visit
// image has width, height, 1d array of bytes
// Returns byte indices of 8 neighbours.
//func neighbours (image []int, i int, start int) []int {
//      0   1   2
//      7       3
//      6   5   4
//func neighbours(img *image.NRGBA, p Point, w, x, y int) {
//	mask := make([]Point, 8);
//  if x == 0 {
//	  // left edge
//    mask[0] = -1;
//	mask[6] = -1;
//	mask[7] = -1;
//  }
//	if x == w-1 {
//    mask[2] = -1;
//	mask[3] = -1;
//	mask[4] = -1;
//  }
//
//  // hack - vertical edging matters less because
//  // it will get ignored by matching it to the source
//
//  // start is always 0 ?
//  //return offset([
//    mask[0] || i - w - 1,
//    mask[1] || i - w,
//    mask[2] || i - w + 1,
//    mask[3] || i + 1,
//    mask[4] || i + w + 1,
//    mask[5] || i + w,
//    mask[6] || i + w - 1,
//    mask[7] || i - 1
//  ];
//}

// Given an image, a point, and an index into the
// relative neighbours array, return a boolean value
// indicating whether it's in the thing or not;
// also returns the next point to be visited.
/*
func getValue(img *image.NRGBA, p PointT, idx int) (bool, PointT) {
	const threshold int = threshold // tunable
	var p2 PointT
	p2.x = p.x + relNeighbours[idx].x
	p2.y = p.y + relNeighbours[idx].y
	if p2.x < 0 || p2.y < 0 || p2.x >= opts.width || p2.y >= opts.height {
		// off the edge
		return false, nilPoint
	}
	//imgValue := img.At(p2.x, p2.y)
	c := color.NRGBAModel.Convert(img.At(p2.x, p2.y)).(color.NRGBA)
	// c has .R, .G, .B, .A (uint8's)  For now, just use .R
	if int(c.R) < threshold {
		return true, p2
	}
	return false, nilPoint
}
*/
/*
func getValue(img *image.NRGBA, p PointT) bool {
	const threshold int = threshold //  FIXME tunable
	if p.x < 0 || p.y < 0 || p.x >= opts.width || p.y >= opts.height {
		// off the edge
		return false
	}
	c := color.NRGBAModel.Convert(img.At(p.x, p.y)).(color.NRGBA) // FIXME this code repeated below
	//imgValue := img.At(p2.x, p2.y)
	// c has .R, .G, .B, .A (uint8's)  For now, just use .R
	if int(c.R) < threshold {
		return true
	}
	return false
}

// 6  5  4
// 7     3
// 0  1  2   +ve y is down
var Steps = []PointT{ // 'relative points'
	{-1, +1}, // 0
	{0, +1},  // 1
	{+1, +1}, // 2
	{+1, 0},  // 3
	{+1, -1}, // 4
	{0, -1},  // 5
	{-1, -1}, // 6
	{-1, 0},  // 7
}

func (p PointT) Step(d int) PointT {
	return p.Plus(Steps[d])
}
func (p PointT) Backtrack(d int) PointT {
	return p.Minus(Steps[d])
}
*/
// Backtrack from p1 by going one step in the opposite direction to dir
//func backtrack(p1 PointT, dir int) PointT {
//}
/*
func traceContour(img *image.NRGBA, start PointT, dir int) ContourT {
	// dir is the direction we came in from
	contour := make(ContourT, 1, 10) // arbitrary starting size
	contour[0] = start
	//var direction = 3
	var p = start.Backtrack(dir)
	backtrackDir := (dir + 4) % 8
	fmt.Printf("tC: starting at %v, dir=%v, p=%v, bdir=%v, \n", start, dir, p, backtrackDir)
	for true {
		//var n = neighbours(imageData, p, 0);
		// find the first neighbour starting from
		// the direction we came from
		// i.e. backtrack
		//var offset = direction - 3 + 8
		   directions:
		     0  1  2
		     7     3
		     6  5  4
		   start indexes:
		     5  6  7
		     4     0
		     3  2  1
		//direction = -1;
		var nextP PointT
		var value bool // TODO rename
		for i := 0; i < 8; i++ {
			// Loop round neighbours, starting one clockwise from backtrack direction
			// idx is index into neighbours array
			//idx := (i1 + offset) % 8
			neighbour := (i + backtrackDir) % 8
			//if(imageData.data[n[idx] * 4] > 0) {
			value = getValue(img, neighbour)
			//if img.At(n[idx].x, n[idx].y) > 0 {
			if value {
				direction = idx
				break // from for i1...
			}
		}
		//p = n[direction];
		p = nextP
		if p.Equal(start) { //  || !p) {
			break
		} else {
			contour = append(contour, p)
		}
	}
	fmt.Printf("tC: returning contour %v\n", contour)
	return contour
}
*/
//
//func offset (array, by) {
//	return array.map( function (_v, i) { return array[(i + by) % array.length]; }
//}

// Given an image file name, return contours as a slice of slices of points.
/*
func wrongcontourFinder(img *image.NRGBA) []ContourT {
	var contours = make([]ContourT, 1, 10)
	var seen [][]bool
	var skipping = false
	seen = make([][]bool, opts.width)
	for j := 0; j < opts.height; j++ {
		seen[j] = make([]bool, opts.height)
	}
	for x := 0; x < opts.width; x++ {
		fmt.Printf("cF: starting row %v\n", x)
		for y := 0; y < opts.width; y++ {
			//if(imageData.data[i * 4] > threshold) {
			// if img.At(x, y) > threshold { // FIXME we seem to be doing this use of threshold twice
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			fmt.Printf("cF: c at %v,%v is %#v\n", x, y, c)
			if int(c.R) < threshold {
				if seen[x][y] || skipping {
					fmt.Printf("cF: skipping %v,%v\n", x, y)
					skipping = true
				} else {
					contour := traceContour(img, PointT{x, y}, 1) // dir=1 because we're scanning y after x
					contours = append(contours, contour)
					// FIXME I think contours can overlap -- or maybe only when we're doing contours at different levels
					// MAYBE blank out pixels in image rather than having a 2D array of flags
					for _, p := range contour {
						seen[p.x][p.y] = true
						// FIXME need to hollow out the shape too
					}
				}
			} else {
				skipping = false
			}
		}
	}
	return contours
}
*/

func offset(array []int, by int) []int {
	array2 := make([]int, len(array))
	//  array.map( (_v, i) => array[(i + by) % array.length])
	for i := range array {
		array2[i] = array[(i+by)%len(array)]
	}
	return array2
}

// List of neighbours to visit, in clockwise order from
func Xneighbours(image *image.NRGBA, width int, i, start int) [8]int {
	w := width
	/* ?
	// convert to x,y
	x := i % w
	y := i - (w * x)
	list := [8]int	// initialised to 0 == black
	offEdge := [8]int	// initialised to false
	if x < 0 || y < 0 || x >= width || y >= width {
	*/
	var mask [8]int
	if (i % w) == 0 {
		mask[0] = -1 // left edge
		mask[6] = -1
		mask[7] = -1
	} else {
		mask[0] = i - w - 1 // goes negative for first row
		mask[6] = i + w - 1
		mask[7] = i - 1
	}
	if ((i + 1) % w) == 0 { // right edge
		mask[2] = -1
		mask[3] = -1
		mask[4] = -1
	} else {
		mask[2] = i - w + 1
		mask[3] = i + 1
		mask[4] = i + w + 1
	}
	mask[1] = i - w
	mask[5] = i + w
	// hack - vertical edging matters less because
	// it will get ignored by matching it to the source
	// +-------+-------+-------+
	// | i-w-1 |  i-w  | i-w+1 |
	// +-------+-------+-------+
	// |  i-1  |   i   |  i+1  |
	// +-------+-------+-------+
	// | i+w-1 |  i+w  | i+w+1 |
	// +-------+-------+-------+
	//return offset([]int{
	//	mask[0] || i-w-1,
	//	mask[1] || i-w,
	//	mask[2] || i-w+1,
	//	mask[3] || i+1,
	//	mask[4] || i+w+1,
	//	mask[5] || i+w,
	//	mask[6] || i+w-1,
	//	mask[7] || i-1,
	//}, start)
	return mask
}

// Get the pixel value (0..255) at the given offset into the image
func getPix(imageData *image.NRGBA, i int) int {
	return int(imageData.Pix[i*4]) // For now, just get the red part of the RGBA
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
func neighboursWithin(imageData *image.NRGBA, width int, p int) ([8]int, [8]bool) {
	var neighbours [8]int
	var within [8]bool // initialised to false
	for i := 0; i < 8; i++ {
		neighbour := p + neighbourOffset[i].x + neighbourOffset[i].y*width
		// convert to x,y
		x := neighbour % width
		y := neighbour / width // integer division
		// check we're not off the edge of the image
		if x >= 0 && x < width && y >= 0 && y < width {
			within[i] = (getPix(imageData, neighbour) < threshold)
		}
		neighbours[i] = neighbour
	}
	fmt.Printf("nW returning %v\n", within)
	return neighbours, within
}

func traceContour(imageData *image.NRGBA, width int, imageLen int, i int) ContourT {
	start := i
	contour := make(ContourT, 1, 10)
	contour[0] = start
	var direction int = 3
	p := start
	fmt.Printf("\ntC: start=%v\n", start)
	for true {
		//n := Xneighbours(imageData, width, p, 0)
		neighbours, withins := neighboursWithin(imageData, width, p)
		// find the first neighbour starting from
		// the direction we came from
		var offset int = direction - 3 + 8
		/*
		   directions:
		     0   1   2
		     7       3
		     6   5   4

		   start indexes:
		     5  6   7
		     4      0
		     3  2   1
		*/
		//direction = -1
		nextP := -1
		//idx := 0
		for i := 0; i < 8; i++ {
			idx := (i + offset) % 8
			//neighbour := neighbours[idx] // n[idx]
			within := withins[idx]
			//fmt.Printf("tC loop: p=%v n=%v  offset=%v idx=%v n[idx]=%v\n", p, n, offset, idx, n[idx])
			fmt.Printf("tC loop: p=%v  offset=%v idx=%v ns=%v ws=%v\n", p, offset, idx, neighbours, withins)
			//if neighbour >= 0 && neighbour < imageLen && imageData.Pix[neighbour*4] < threshold { // > 0 { // > threshold?
			if within {
				direction = idx
				nextP = neighbours[idx]
				//fmt.Printf("tC: breaking with direction=%v nextP=%v\n", direction, nextP)
				fmt.Printf("tC: breaking with nextP=%v\n", nextP)
				break
			}
		}
		if nextP == -1 {
			fmt.Printf("tC: !!!!!!!! finished loop without breaking\n")
		}
		fmt.Printf("tC: old p=%v  nextP=%v\n", p, nextP)
		p = nextP                  // neighbours[direction]
		if p == start || p == -1 { //?? || (p != 0) {	//!p {
			break
		} else {
			contour = append(contour, p)
		}
	}
	return contour
}

func contourFinder(imageData *image.NRGBA, width, height int) []ContourT {
	var contours = make([]ContourT, 1, 10)
	var imageLen = width * height
	seen := make([]bool, imageLen)
	var skipping = false
	for i := 0; i < imageLen; i++ {
		//if(imageData.data[i * 4] > threshold) {
		// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
		if imageData.Pix[i*4] < threshold { // just get R
			if seen[i] || skipping {
				skipping = true
			} else {
				var contour = traceContour(imageData, width, imageLen, i)
				contours = append(contours, contour)
				// this could be a _lot_ more efficient
				//contour.forEach(c => {
				fmt.Printf("cF: got contour %v\n", contour)
				for _, c := range contour {
					seen[c] = true
				}
			}
		} else {
			skipping = false
		}
	}
	return contours
}

func parseArgs(args []string) {

	// Use a flagset because it allows passing command line arguments as a slice of strings
	pf := pflag.NewFlagSet("contours", pflag.ExitOnError)

	//pf.Float64VarP(&opts.angleDegrees, "angle", "a", 90, "Angle between gradient and contour: -90..90 degrees.")
	//pf.IntVarP(&opts.cellSize, "cellsize", "c", 3, "Cellsize: 3, 5, or 7.")
	//pf.BoolVarP(&opts.darker, "darker", "d", false, "Darker shading (more contours).  Implies --multi. (default false)")
	//pf.Float64Var(&opts.margin, "margin", 15, "Minimum margin (in mm).")
	//pf.BoolVarP(&opts.multiContour, "multi", "m", false, "Use multiple contours in each cell (depending on cell size). (default false)")
	//pf.BoolVarP(&opts.gradient, "gradient", "g", false, "Set number of contours depending on gradient (rather than darkness). (default false)")
	//pf.IntVarP(&opts.overlap, "overlap", "o", 0, "Overlap between cells: 0..cellsize-1. (default 0)")
	//pf.StringVar(&opts.paper, "paper", "A4L", "Paper size and orientation.  A4L | A4P | A3L | A3P.")
	//pf.IntVarP(&opts.penCount, "pens", "n", 1, "Number of pens to use (for shading, e.g. grey and black would be 2). (default 0)")
	//pf.BoolVarP(&opts.PNGoutput, "png", "p", false, "Create PNG output file. (default false)")
	//pf.BoolVarP(&opts.SVGoutput, "svg", "s", false, "Create SVG output file. (default true if no --png)")
	//pf.BoolVarP(&opts.varLen, "varlen", "v", false, "Use variable-length lines to represent shading. (default false)")

	pf.SortFlags = false
	if args == nil {
		pf.Parse(os.Args[1:]) // don't pass program name
	} else {
		pf.Parse(args) // args passed as a string (for testing)
	}

	// Validate
	//ok := true
	//if opts.cellSize != 3 && opts.cellSize != 5 && opts.cellSize != 7 {
	//		fmt.Printf("Invalid cellSize '%v' -- should be 3 or 7\n", opts.cellSize)
	//		ok = false
	//} else if opts.overlap < 0 || opts.overlap > (opts.cellSize-1) {
	//	fmt.Printf("Invalid overlap `%v` -- should be between 0 and %d\n", opts.overlap, opts.cellSize-1)
	//	ok = false
	//}
	//if opts.angleDegrees < -90 || opts.angleDegrees > 90 {
	//	fmt.Printf("Invalid angle `%v` -- should be between -90 and 90\n", opts.angleDegrees)
	//	ok = false
	//}
	//if opts.margin < 0 || opts.margin > 200 {
	//	fmt.Printf("Invalid margin `%v` -- should be between 0 and 200 mm\n", opts.margin)
	//	ok = false
	//}
	//if opts.penCount < 1 || opts.penCount > 9 {
	//	fmt.Printf("Invalid number of pens `%v` -- should be between 1 and 9 mm\n", opts.penCount)
	//	ok = false
	//}
	//opts.paper = strings.ToUpper(opts.paper)
	//if _, paperOk := paperSizes[opts.paper]; !paperOk {
	//	fmt.Printf("Invalid paper size '%v' -- should be one of %v\n", opts.paper, paperSizes) // TODO better error message
	//	ok = false
	//}
	//if pf.NArg() < 1 {
	//	fmt.Println("No input file name given")
	//	ok = false
	//}
	//if !ok {
	//	os.Exit(1)
	//}
	opts.infile = pf.Arg(0)

	// Make sure there's at least one output file
	//if !opts.PNGoutput && !opts.SVGoutput {
	//	opts.SVGoutput = true
	//}
	// --darker implies --multi
	//if opts.darker {
	//	opts.multiContour = true
	//}

	//ext := filepath.Ext(opts.infile)
	//multiFlag := ""
	//if opts.multiContour {
	//	multiFlag = "M"
	//}
	//varlenFlag := ""
	//if opts.varLen {
	//	varlenFlag = "V"
	//}
	//darkerFlag := ""
	//if opts.darker {
	//	darkerFlag = "D"
	//}
	//pensFlag := ""
	//if opts.penCount > 1 {
	//	pensFlag = fmt.Sprint(opts.penCount)
	//}
	//optString := fmt.Sprintf("-c%do%da%.0f%s%s%s%s-m%g%s",
	//	opts.cellSize, opts.overlap, opts.angleDegrees, multiFlag, varlenFlag, darkerFlag, pensFlag, opts.margin, opts.paper)
	//png.filename = strings.TrimSuffix(opts.infile, ext) + optString + ".png"
	//svgF.filename = strings.TrimSuffix(opts.infile, ext) + optString + ".svg"

}

func main() {
	parseArgs(nil)
	fmt.Printf("mncontours: processing '%s'\n", opts.infile)
	img, width, height, err := loadImage(opts.infile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	opts.width = width
	opts.height = height
	fmt.Printf("options: %#v\n", opts)
	contours := contourFinder(img, width, height)
	fmt.Printf("contours: %#v\n", contours)
}
