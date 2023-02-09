package test

import (
	"fmt"
	"qsr-mock-server/utils"
	"strings"
	"testing"
	"time"
)

func TestGetAllIndex(t *testing.T) {
	srcStr := "qsrhsjfkhsqsrkchskjchsqsrlsdfjs"
	tarStr := "qsr"
	indexs := []int{}
	utils.GetAllIndex(srcStr, tarStr, 0, &indexs)
	byteStr := []byte(srcStr)
	t.Log(indexs)
	t.Log(string(byteStr[0]))
	t.Log(byteStr[10]) // 直接获取的是 ASCII 码
	t.Log(byteStr[22])
}

func TestJsonStrContain(t *testing.T) {
	//jsonStr1 := "{\"json\": \"post\",\"message\": \"1234\"}"
	//jsonStr2 := "{\"json\": \"post\",\"message\"             :               \"1234\"}"
	jsonStr3 := "{\"json\": \"post\",\"message\"             :    \n\t           \"1234\"}"
	contain := utils.JsonStrContain(jsonStr3, "message", "1234")
	t.Log(contain)
}

// 多层循环，break 只能破坏最内层
func TestOne(t *testing.T) {
	n := 20
	for {
		i := 0
		for {
			if i == 10 {
				break
			}
			i += 1
			fmt.Println("1111111")
		}
		n -= 1
		time.Sleep(time.Second)
		fmt.Println(n)
	}
}
func TestOne1(t *testing.T) {
	str := "http://sakjda"
	strs := strings.Split(str, "?")
	for _, s := range strs {
		t.Log(s)
	}
}
