package ginutil

import (
	"fmt"
	"testing"
)

func TestGetPageNation(t *testing.T) {
	n := GetPageNation(nil, 10)
	fmt.Println(n)
}
