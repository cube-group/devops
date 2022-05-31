package slice

import (
	"fmt"
	"testing"
)

func TestSliceStringDiff(t *testing.T) {
	s1 := []string{"b", "c", "e"}
	s2 := []string{"x", "e"}
	s3 := []string{"b", "c"}
	if SliceStringEqual(s3, SliceStringDiff(s1, s2)) {
		fmt.Println("TestSliceStringDiff ok")
	} else {
		t.Error("TestSliceStringDiff Err.")
	}
}

func TestSliceStringMixed(t *testing.T) {
	s1 := []string{"b", "c", "e"}
	s2 := []string{"x", "e"}
	s3 := []string{"e"}
	if SliceStringEqual(s3, SliceStringMixed(s1, s2)) {
		fmt.Println("TestSliceStringMixed ok")
	} else {
		t.Error("TestSliceStringMixed Err.")
	}
}

func TestSliceStringUnion(t *testing.T) {
	s1 := []string{"b", "c", "e"}
	s2 := []string{"x", "e"}
	s3 := []string{"b", "c", "e", "x"}
	if SliceStringEqual(s3, SliceStringUnion(s1, s2)) {
		fmt.Println("TestSliceStringUnion ok")
	} else {
		t.Error("TestSliceStringUnion Err.")
	}
}
