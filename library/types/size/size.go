package size

import (
	"fmt"
	"strconv"
)

func SizeFormat(size interface{}) string {
	res, err := strconv.ParseFloat(fmt.Sprintf("%v", size), 0)
	if err != nil {
		fmt.Println(err)
		return "0 Bytes"
	}
	newSize := float64(res)

	sizeConfig := []string{
		0: "Bytes", 1: "KB", 2: "MB", 3: "GB", 4: "TB",
	}

	var resIndex int

	for newSize >= 1024 {
		newSize = float64(newSize / 1024)
		resIndex++
		if resIndex == 4 {
			break
		}
	}
	return fmt.Sprintf("%.1f %v", newSize, sizeConfig[resIndex])
}

func SizeFormatToMBWithoutUnit(size interface{}) string {
	res, err := strconv.ParseFloat(fmt.Sprintf("%v", size), 0)
	if err != nil {
		fmt.Println(err)
		return "0"
	}
	newSize := float64(res)
	newSize = newSize / 1024 / 1024
	return fmt.Sprintf("%.2f", newSize)
}
