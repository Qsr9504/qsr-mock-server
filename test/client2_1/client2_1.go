package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"qsr-mock-server/utils"
)

// 用户名
var loginName string

// 本机连接
var selfConnect *net.TCPConn

// 读取行文本
var reader = bufio.NewReader(os.Stdin)

// 建立连接
func connect(addr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr) // 使用tcp
	con, err := net.DialTCP("tcp", nil, tcpAddr)  // 拨号：主动向server建立连接
	selfConnect = con
	if err != nil {
		fmt.Println("连接服务器失败")
		os.Exit(1)
	}
	go msgSender()
	go msgReceiver()
}

// 消息接收器
func msgReceiver() {
	fmt.Println("msgReceiver", utils.GetGid())
	buff := make([]byte, 2048)
	for {
		len, _ := selfConnect.Read(buff) // 从建立连接的缓冲区读消息
		fmt.Println(string(buff[:len]))
	}
}

// 消息发送器
func msgSender() {
	fmt.Println("msgSender", utils.GetGid())
	for {
		bMsg, _, _ := reader.ReadLine()
		bMsg = []byte(loginName + " : " + string(bMsg))
		selfConnect.Write(bMsg) // 发消息
	}
}

// 初始化
func initGroupChatClient() {
	fmt.Println("请问您怎么称呼？")
	bName, _, _ := reader.ReadLine()
	loginName = string(bName)
	connect("127.0.0.1:8789")
	select {}
}

func main() {
	initGroupChatClient()
}
