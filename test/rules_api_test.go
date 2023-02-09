package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//   /mockRule/addRule
func TestOne_01(t *testing.T) {
	url := "http://localhost:8787/mockRule/addRule"
	method := "POST"
	payload := strings.NewReader(`{
		"OwnerName":"demoName",
		"ProxyMethod":"ANY",
		"ProxyHost":"localhost:8081",
		"ProxyUrl":"/ping",
		"ProxyConditionJson":"",
		"ProxyConditions":[
			{
				"InterceptSite":0,
				"IsContain":true,
				"InterceptKey":"qsr_key",
				"InterceptValue":"123"
			},
			{
				"InterceptSite":1,
				"IsContain":false,
				"InterceptKey":"qsr_key",
				"InterceptValue":"123"
			}
		],
		"CustomizedResp":"{\"message\":\"ok\"}"
	}`)

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	readAll, _ := ioutil.ReadAll(res.Body)
	t.Log(string(readAll))
}

//  /mockRule/deleteRulesById
func TestOne_02(t *testing.T) {
	url := "http://localhost:8787/mockRule/deleteRulesById"
	method := "POST"

	payload := strings.NewReader(`{"id":7}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}

//  /mockRule/updateRulesById
func TestOne_03(t *testing.T) {
	url := "http://localhost:8787/mockRule/updateRulesById"
	method := "POST"

	payload := strings.NewReader(`{
    "id":3,
    "OwnerName":"demoName",
    "ProxyMethod":"ANY",
    "ProxyHost":"localhost:8082",
    "ProxyUrl":"/ping",
    "ProxyConditionJson":"",
    "ProxyConditions":[
        {
            "InterceptSite":0,
            "IsContain":true,
            "InterceptKey":"0",
            "InterceptValue":"0"
        },
        {
            "InterceptSite":0,
            "IsContain":true,
            "InterceptKey":"1",
            "InterceptValue":"1"
        }
    ],
    "CustomizedResp":"{\"message\":\"ok\"}"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}

// /mockRule/getAllRules
func TestOne_04(t *testing.T) {
	url := "http://localhost:8787/mockRule/getAllRules"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}

// /mockRule/getRuleDetailById
func TestOne_05(t *testing.T) {
	url := "http://localhost:8787/mockRule/getRuleDetailById"
	method := "GET"

	payload := strings.NewReader(`{"id":3}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}

// /mockRule/getRuleByCons
func TestOne_06(t *testing.T) {
	url := "http://localhost:8787/mockRule/getRuleByCons"
	method := "GET"

	payload := strings.NewReader(`{"ProxyMethod":"POST"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}

// 一个成功的路由代理测试
func TestOne_07(t *testing.T) {

	url := "http://localhost:8787/mock/localhost:8081/ping?hello=123&nihao=789&qsr=qqqqqqq"
	method := "GET"

	payload := strings.NewReader(`{"qsr_key":"123"}`) // 修改为 qsr_key 会报错 TODO：待排查

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("qsr_key", "123")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ownerName", "demoName") // 添加用户信息

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(body))
}
