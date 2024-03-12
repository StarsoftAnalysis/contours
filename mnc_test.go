package main

import (
	"fmt"
	"image"
	"testing"
)

func TestTraceContour(t *testing.T) {
	fmt.Println("TestTraceContour")
	var infile string
	var img *image.NRGBA
	var width int
	var err error
	var start int
	var wanted, got ContourT

	infile = "test0.png"
	img, width, _, err = loadImage(infile)
	if err != nil {
		t.Errorf("Input file %s not found\n", infile)
	}
	start = 5
	wanted = ContourT{5, 6, 10, 9}
	got = traceContour(img, width, start)
	if !got.Equals(wanted) {
		t.Errorf("Wrong result for test0.png/%v (wanted=%v  got %v)\n", start, wanted, got)
	}

	infile = "test1.png"
	img, width, _, err = loadImage(infile)
	if err != nil {
		t.Errorf("Input file %s not found\n", infile)
	}
	start = 6
	wanted = ContourT{6, 7, 8, 12, 16, 11}
	got = traceContour(img, width, start)
	//fmt.Printf("result for test1.png/%v (wanted=%v  got %v)\n", start, wanted, got)
	if !got.Equals(wanted) {
		t.Errorf("Wrong result for test1.png/%v (wanted=%v  got %v)\n", start, wanted, got)
	}
}
