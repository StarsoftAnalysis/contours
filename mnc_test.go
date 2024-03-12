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
		{"examples/test4.png", ContourT{1, 8, 14, 19, 12, 6}, 1},
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
		{"examples/test2.png", []ContourT{{0, 1, 4}, {11, 15, 14}}, 2},
		{"examples/test3.png", []ContourT{{1, 2, 10, 17, 16, 8}, {4, 5, 6, 7, 15, 23, 31, 39, 47, 55, 63, 62, 61, 60, 59, 58, 57, 56, 48, 40, 32, 33, 34, 27, 20, 12}, {37, 45, 44}}, 3},
		{"examples/test4.png", []ContourT{{1, 8, 14, 19, 12, 6}, {4, 5, 11, 10}, {22}}, 3},
		{"examples/example.jpg", nil, 13},
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
