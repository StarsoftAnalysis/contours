package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTraceContour(t *testing.T) {
	fmt.Println("TestTraceContour")

	type testdataT struct {
		infile  string
		contour ContourT
		start   int
	}
	testdata := []testdataT{
		{"examples/test0.png", ContourT{5, 6, 10, 9}, 5},
		{"examples/test1.png", ContourT{6, 7, 8, 12, 16, 11}, 6},
		{"examples/test4.png", ContourT{5, 6, 10, 9}, 2},
	}

	for _, td := range testdata {
		img, width, height, err := loadImage(td.infile)
		if err != nil {
			t.Errorf("Input file %s not found\n", td.infile)
		}
		got := traceContour(img, width, height, td.start)
		if !reflect.DeepEqual(got, td.contour) {
			t.Errorf("Wrong result for %s/%v (wanted=%v  got %v)\n", td.infile, td.start, td.contour, got)
		}
	}
}

func TestContourFinder(t *testing.T) {
	fmt.Println("TestContourFinder")
	type testdataT struct {
		infile   string
		contours []ContourT
		count    int
	}
	testdata := []testdataT{
		{"examples/test0.png", []ContourT{{5, 6, 10, 9}}, 1},
		{"examples/test1.png", []ContourT{{6, 7, 8, 12, 16, 11}}, 1},
		//{"examples/test2.png", []ContourT{{0, 1, 4}, {11, 15, 14}}, 2},
		//{"examples/test3.png", []ContourT{{1, 2, 10, 17, 16, 23, 32, 33, 34, 27, 20, 12, 4, 5, 6, 7, 8}, {15, 6, 5, 12, 20, 27, 34, 33, 40, 48, 57, 58, 59, 60, 61, 62, 55, 47, 39, 31, 23}, {37, 45, 44}}, 3},
		{"examples/test4.png", []ContourT{{5, 6, 10, 9}}, 1},
		//{"examples/example.jpg", nil, 13},
	}
	for _, td := range testdata {
		img, width, height, err := loadImage(td.infile)
		if err != nil {
			t.Errorf("Input file %s not found\n", td.infile)
		}
		got := contourFinder(img, width, height)
		if td.contours == nil {
			// just count the contours
			if len(got) != td.count {
				t.Errorf("Wrong result for %s (wanted length %v  got %v)\n", td.infile, td.count, len(got))
			}
		} else {
			if !reflect.DeepEqual(got, td.contours) {
				t.Errorf("Wrong result for %s (wanted=%v  got %v)\n", td.infile, td.contours, got)
			}
		}
	}
}
