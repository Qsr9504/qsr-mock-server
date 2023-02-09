package main

import (
	"fmt"
	"io"
	"net"
)

func serveForClient(conn net.Conn) {
	// 我们这里循环的 接收客户端发来的消息
	defer conn.Close()
	for {
		// 创建一个新的切片
		buf := make([]byte, 10240)
		// conn.Read()
		// 1. 等待客户端通过conn发送消息
		// 2. 如果客户端没有 write 发送，那么这个协程整个就会阻塞在这里
		fmt.Printf("服务端陷入等待 %s 来消息\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取客户端发来的消息失败:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}

func main() {
	fmt.Println("服务端开始监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8788")
	if err != nil {
		fmt.Println("开始监听出错")
		return
	}
	fmt.Println(listen)
	defer listen.Close()
	// 循环等待链接
	for {
		// 等待客户端来连接
		fmt.Println("等待客户端连接")
		con, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept 等待监听出错")
			if err == io.EOF {
				fmt.Println("客户端已经退出")
			}
			return
		} else {
			fmt.Printf("Accept() suc con=%v\n", con)
		}
		// 这里准备启动一个协程，为客户端服务
		go serveForClient(con)
	}
}
