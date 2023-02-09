// client.go
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

func main() {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := "ws://localhost:8787" + "/spider-service/v1/channels"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)
	var sendStr string
	sendStr = "{\"hi\":{\"id\":\"111084\",\"ver\":\"0.18\",\"ua\":\"TinodeWeb/0.18 (Chrome/109.0; MacIntel); dovejs/0.18\",\"lang\":\"zh-CN\",\"platf\":\"web\",\"appver\":\"5.6.0\",\"userId\":\"10119168\",\"domain\":0},\"login\":{\"id\":\"111084\",\"scheme\":\"kefu\",\"secret\":\"YzAzZjBmOWZjMjkxNGIxNGJlODBiYmY3NDdjMDNjZDI=\",\"domain\":0}}"
	conn.WriteMessage(websocket.TextMessage, []byte(sendStr))
	time.Sleep(time.Second * 10)
	sendStr = "{\"get\":{\"id\":\"111103\",\"topic\":\"grpHbCpv2zU7oM\",\"what\":\"data\",\"data\":{\"limit\":20,\"gettype\":1,\"batchres\":true},\"isHistory\":true,\"domain\":0}}"
	conn.WriteMessage(websocket.TextMessage, []byte(sendStr))
	select {}

	// 无限循环使用select来通过通道监听事件
	//for {
	//	select {
	//	case <-time.After(time.Duration(1) * time.Millisecond * 1111):
	//		//conn.WriteMessage()每秒钟写一条消息
	//		err := conn.WriteMessage(websocket.TextMessage, []byte("qsr!"))
	//		if err != nil {
	//			log.Println("Error during writing to websocket:", err)
	//			return
	//		}
	//	//如果激活了中断信号，则所有未决的连接都将关闭
	//	case <-interrupt:
	//		// We received a SIGINT (Ctrl + C). Terminate gracefully...
	//		log.Println("Received SIGINT interrupt signal. Closing all pending connections")
	//
	//		// Close our websocket connection
	//		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			log.Println("Error during closing websocket:", err)
	//			return
	//		}
	//
	//		select {
	//		// 如果receiveHandler通道退出，则通道'done'将关闭
	//		case <-done:
	//			log.Println("Receiver Channel Closed! Exiting....")
	//		//如果'done'通道未关闭，则在1秒钟后会有超时，因此程序将在1秒钟超时后退出
	//		case <-time.After(time.Duration(1) * time.Second):
	//			log.Println("Timeout in closing receiving channel. Exiting....")
	//		}
	//		return
	//	}
	//}
}
