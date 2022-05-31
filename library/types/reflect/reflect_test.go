package reflect

import (
	"fmt"
	"testing"
	"time"
)

type TestStructDemo struct {
	Name string    `json:"name"`
	Age  int       `json:"age"`
	T    time.Time `json:"time"`
}

func TestGetStructMap(t *testing.T) {
	res, err := GetStructMap(TestStructDemo{"linyang123", 32, time.Now()})
	if err != nil {
		t.Errorf("TestGetStructMap Err: %s", err.Error())
	} else {
		fmt.Println(res)
	}
}

func TestGetStructFields(t *testing.T) {
	res, err := GetStructFields(TestStructDemo{})
	if err != nil {
		t.Errorf("TestGetStructFields Err: %s", err.Error())
	} else {
		t.Log(res)
	}
}

func TestGetStructTags(t *testing.T) {
	res, err := GetStructTags(TestStructDemo{})
	if err != nil {
		t.Errorf("TestGetStructTags Err: %s", err.Error())
	} else {
		t.Log(res)
	}
}

func TestNoStruct(t *testing.T) {
	res, err := GetStructTags("a")
	if err != nil {
		t.Errorf("TestGetStructTags Err: %s", err.Error())
	} else {
		t.Log(res)
	}
}
