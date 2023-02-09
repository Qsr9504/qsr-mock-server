package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8788")
	if err != nil {
		fmt.Println(" client dial err = ", err)
		return
	}
	defer conn.Close()
	// 功能1：客户端可以发送单行数据，然后就退出
	reader := bufio.NewReader(os.Stdin) // 终端输入的内容

	for {
		// 从终端读取一行用户输入，并且发送给服务器
		lineStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("readString")
		}

		// 如果用户输入的是exit，就退出
		if strings.Trim(lineStr, "\n\r") == "exit" {
			fmt.Println("客户端退出")
			break
		}

		// 再将lineStr进行发送给服务器
		n, err := conn.Write([]byte(lineStr + "\n"))
		if err != nil {
			fmt.Println("客户端阀发送给服务器报错")
		}
		fmt.Println("客户端发送了 %d 个字节的数据", n)
	}
}
