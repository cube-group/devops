package slice

import "testing"

func TestInArray(t *testing.T) {
	source := []string{"a", "b"}
	if !InArray("a", source) {
		t.Fatal("no")
	}
}
