// Author: chenqionghe
// Time: 2018-10
// 调试打印

package ginutil

import (
	"fmt"
	"strings"
)

//调试打印
func Dump(v ...interface{}) {

	fmt.Print(strings.Repeat("*", 49), " begin ", strings.Repeat("*", 49))
	fmt.Println(strings.Repeat("\r\n", 3))
	for _, i := range v {
		fmt.Println(i)
	}
	fmt.Println(strings.Repeat("\r\n", 3))
	fmt.Print(strings.Repeat("*", 50), " end ", strings.Repeat("*", 50))
	fmt.Println()
}
