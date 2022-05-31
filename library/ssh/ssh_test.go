package ssh

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	obj, err := NewByKey("39.106.107.76", "root", "/Users/chenqionghe/web/offcn/base/corecd-v2/conf/sshkey", 22)
	_, _, err = obj.Run(`

abc
echo 'a'
sleep 1

echo 'b'
sleep 1

echo 'c'

sleep 1

echo 'd'
exit 0
`)
	fmt.Println(err)
}
