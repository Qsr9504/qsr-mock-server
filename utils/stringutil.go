package utils

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"
	"strings"
)

// JsonStringToMap JSON字符串转map
func JsonStringToMap(content string) (map[string]interface{}, error) {
	var resMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &resMap)
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

// JsonStrContain 判断一个 jsonStr JSON字符串 中是否包含 key value 这一对参数，false 不包含
func JsonStrContain(jsonStr, key, value string) (isContain bool) {
	stringMap, err := JsonStringToMap(jsonStr)
	if err != nil {
		return false
	}
	for k, v := range stringMap {
		if k == key && v == value {
			return true
		}
	}
	return false
}

func GetGid() (gid uint64) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

/*

// JsonStrContain 判断一个 jsonStr JSON字符串 中是否包含 key value 这一对参数，false 不包含
func JsonStrContain(jsonStr, key, value string) (isContain bool) {
	isContain = false
	//"json":"post"
	bytesJson := []byte(jsonStr)
	keyLen := len(key)
	valueLen := len(value)
	jsonStrLen := len(jsonStr)
	keyIndexs := []int{}
	valueIndexs := []int{}
	GetAllIndex(jsonStr, key, 0, &keyIndexs)
	GetAllIndex(jsonStr, value, 0, &valueIndexs)
	fmt.Println("keyIndexs", keyIndexs)
	fmt.Println("valueIndexs", valueIndexs)
	if len(keyIndexs) == 0 || len(valueIndexs) == 0 {
		fmt.Println("集合有一个为空，必然不存在")
		return
	}
	for _, keyIndex := range keyIndexs {
		for _, valueIndex := range valueIndexs {
			if valueIndex <= 4 || valueIndex <= keyIndex+keyLen || keyIndex+keyLen+2 > jsonStrLen {
				fmt.Println("不符合基本条件")
				return
			}
			if keyIndex+keyLen+1+valueLen+1 > jsonStrLen || string(bytesJson[keyIndex+keyLen]) != "\"" {
				// 如果 key 后边紧跟着的 不是 双引号 \"
				// keyIndex+keyLen+1+valueLen+1 ，也就是说 ：keyIndex+keyLen+1 这个是:号，最极端情况，value只有1位，那就是加上一个valueLen的长度，不能超过字符串总长度，
				// 但是，等号左边是索引，右边是长度，因此，需要再+1 ，如果 keyIndex+keyLen+1+valueLen+1 > jsonStrLen的长度，就说明这不是一个key
				// "qsr":1
				break
			}
			gapKeyCount := 0
			for {
				if keyIndex+keyLen+gapKeyCount > jsonStrLen-1 {
					fmt.Println("key 还没找到冒号，就超长了哇")
					return
				}
				c := string(bytesJson[keyIndex+keyLen+gapKeyCount])

				if c == ":" {
					break
				}
				gapKeyCount += 1
			}

			gapValueCount := 0
			for {
				if valueIndex-1-gapValueCount < 0 {
					break
				}
				c := string(bytesJson[valueIndex-1-gapValueCount])
				if c == ":" {
					break
				}

				gapValueCount += 1
			}

			// TODO: 这里有个bug，如果在 值 在最后 ，考虑 {[[["qsr","123"]]]}的情况
			//// 如果 值 的最后一个有效字符的后边一个字符 不是整个json的最后一个   有效字符 "qsr":"123" ，有效字符就是3
			//if valueIndex+valueLen != jsonStrLen-1 {
			//	c := string(bytesJson[valueIndex+valueLen])
			//	// 如果 值 的最后一个有效字符的后边一个 是 整个json的倒数第二个
			//	if valueIndex+valueLen == jsonStrLen-2 {
			//		if c != "\"" {
			//			// 那就必然是双引号，不是就是错误的
			//			return
			//		}
			//	} else {
			//		// 如果不是 倒数第二个
			//		if c == "\"" {
			//
			//		}
			//	}
			//} else {
			//	// 如果 值 的最后一个有效字符，的后边一个字符，是整个json 的最后一个，就符合条件
			//
			//}

			// 确保 值 的最后一位后边一个是双引号，或者是一个逗号
			if c := string(bytesJson[valueIndex+valueLen]); c != "\"" {
				return
			}
			// 确保 值 的第一位前边一个是双引号，或者是一个冒号
			if c := string(bytesJson[valueIndex-1]); c != "\"" || c != ":" {
				return
			}

			if valueIndex-1-gapValueCount == keyIndex+keyLen+gapKeyCount { // 等号左边，是从value的值找到了 : 的下标，右边是从key的值找到了:的下标
				isContain = true
				fmt.Println("对比上了哇！")
				return
			}
		}
	}
	fmt.Println("json匹配键值对 盼到最后，扔出去")
	return
}

*/

// GetAllIndex srcStr,待检索的母串；tarStr，待检索的子串；startIndex 从第几个位置开始检索， indexs 子串所有出现的位置
func GetAllIndex(srcStr, tarStr string, startIndex int, indexs *[]int) {
	index := strings.Index(srcStr[startIndex:], tarStr)
	if index == -1 {
		return
	}
	*indexs = append(*indexs, startIndex+index)                       // 将获取到的 索引位置放入
	GetAllIndex(srcStr, tarStr, startIndex+index+len(tarStr), indexs) // 递归进行检索
}
