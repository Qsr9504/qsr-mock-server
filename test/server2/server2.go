package main

import (
	"fmt"
	"net"
	"qsr-mock-server/utils"
	"time"
)

// 客户端 map
var clientMap = make(map[string]*net.TCPConn) // 存储当前群聊中所有用户连接信息：key: ip+port, val: 用户连接信息

// 监听请求
func listenClient(ipAndPort string) {
	fmt.Println("listenClient:", utils.GetGid())
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ipAndPort)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for { // 循环接收
		clientConn, _ := tcpListener.AcceptTCP()                 // 监听请求连接
		clientMap[clientConn.RemoteAddr().String()] = clientConn // 将连接添加到 map
		go addReceiver(clientConn)
		fmt.Println("用户 : ", clientConn.RemoteAddr().String(), " 已连接.")
	}
}

// 向连接添加接收器
func addReceiver(newConnect *net.TCPConn) {
	for {
		byteMsg := make([]byte, 2048)
		len, err := newConnect.Read(byteMsg) // 从newConnect中读取信息到缓存中
		if err != nil {
			newConnect.Close()
		}
		fmt.Println(string(byteMsg[:len]))
		msgBroadcast(byteMsg[:len], newConnect.RemoteAddr().String())
	}
}

// 广播给所有 client
func msgBroadcast(byteMsg []byte, key string) {
	for k, con := range clientMap {
		if k != key { // 转发消息给当前群聊中，除自身以外的其他用户
			fmt.Println("RemoteAddr().String():", con.RemoteAddr().String())
			fmt.Println("RemoteAddr().Network():", con.RemoteAddr().Network())
			fmt.Println("LocalAddr().String():", con.LocalAddr().String())
			fmt.Println("LocalAddr().Network():", con.LocalAddr().Network())
			con.Write(byteMsg)
		}
	}
}

// 初始化
func initGroupChatServer() {
	fmt.Println(utils.GetGid())
	fmt.Println("服务已启动...")
	time.Sleep(1 * time.Second)
	fmt.Println("等待客户端请求连接...")
	go listenClient("0.0.0.0:8789")
	select {} // 阻塞协程，等待监听
}

func main() {
	initGroupChatServer()
}
