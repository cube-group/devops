package times

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	fmt.Println(FormatDatetime(time.Now()))
}
