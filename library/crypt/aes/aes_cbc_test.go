package aes

import (
	"fmt"
	"testing"
)

func TestAesCBC_Encrypt(t *testing.T) {
	res, err := AesCBC_Encrypt("1234567890123456", "hello")
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(res)
	}
}

func TestAesCBC_Decrypt(t *testing.T) {
	res, err := AesCBC_Decrypt("1234567890123456", "ObBxtb9plyPvM6ZEdBv6MQ==")
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(res)
	}
}
