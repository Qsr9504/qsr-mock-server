package main

import (
	"fmt"
	"testing"
)

func defer_test() {
	defer func() {
		fmt.Println("123")
	}()
	fmt.Println("111")
	return
}

type name struct {
	id string
}

func Test_main(t *testing.T) {
	//defer_test()
	m := make(map[string]name, 10)
	fmt.Println(m["1"].id)
}
