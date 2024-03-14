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
		{"tests/test0.png", ContourT{5, 6, 10, 9}, 5},
		{"tests/test1.png", ContourT{6, 7, 8, 12, 16, 11}, 6},
		{"tests/test4.png", ContourT{1, 8, 14, 19, 12, 6}, 1},
	}
	for _, td := range testdata {
		fmt.Printf("\t%s\n", td.infile)
		img, width, height, err := loadImage(td.infile)
		if err != nil {
			t.Errorf("Input file %s not found\n", td.infile)
		}
		got := traceContour(img, width, height, 128, td.start, nil)
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
		{"tests/test0.png", []ContourT{{5, 6, 10, 9}}, 1},
		{"tests/test1.png", []ContourT{{6, 7, 8, 12, 16, 11}}, 1},
		{"tests/test2.png", []ContourT{{0, 1, 4}, {11, 15, 14}}, 2},
		{"tests/test3.png", []ContourT{{1, 2, 10, 17, 16, 8}, {4, 5, 6, 7, 15, 23, 31, 39, 47, 55, 63, 62, 61, 60, 59, 58, 57, 56, 48, 40, 32, 33, 34, 27, 20, 12}, {37, 45, 44}}, 3},
		{"tests/test4.png", []ContourT{{1, 8, 14, 19, 12, 6}, {4, 5, 11, 10}, {22}}, 3},
		{"tests/test5.png", []ContourT{{9, 10, 11, 12, 13, 14, 22, 30, 38, 46, 54, 53, 52, 51, 50, 49, 41, 33, 25, 17}, {29, 20, 19, 26, 34, 43, 44, 37}}, 2},
		{"tests/test6.png", []ContourT{{9, 10, 11, 12, 13, 14, 22, 30, 38, 46, 54, 53, 52, 51, 50, 49, 41, 33, 25, 17}}, 1},
		// These two have non-closed thin lines -- the contour loops back to close itself:
		{"tests/test7.png", []ContourT{{6, 7, 8, 13, 18, 17, 21, 27, 28, 27, 26, 21, 16, 17, 13, 7}}, 1},
		{"tests/test8.png", []ContourT{{10, 11, 12, 13, 14, 15, 16, 15, 14, 13, 12, 11}, {28, 37, 46, 55, 64, 55, 46, 37}, {30, 40, 50, 60, 70, 60, 50, 40}}, 3},
		{"tests/example.png", nil, 10},
		{"tests/bottom.png", nil, 156},
	}
	for _, td := range testdata {
		fmt.Printf("\t%s\n", td.infile)
		img, width, height, err := loadImage(td.infile)
		if err != nil {
			t.Errorf("Input file %s not found\n", td.infile)
		}
		got := contourFinder(img, width, height, 128, nil)
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
