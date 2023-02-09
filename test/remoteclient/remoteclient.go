package main

import (
	"bufio"
	"github.com/gorilla/websocket"

	"log"
	"os"
)

type websocketClient struct {
	addr        *string
	path        string
	conn        *websocket.Conn
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

func main() {
	dl := websocket.Dialer{}
	conn, _, err := dl.Dial("mockwss://t1-im-gateway.dewu.net/spider-service/v1/channels", nil)
	if err != nil {
		log.Println("建立链接失败！！")
		return
	}
	// err = conn.WriteMessage(websocket.TextMessage, body)
	go send(conn)
	for {
		m, p, e := conn.ReadMessage()
		if e != nil {
			return
		}
		log.Println(m, string(p))
	}
}

func send(conn *websocket.Conn) {
	for {

		reader := bufio.NewReader(os.Stdin)
		l, _, _ := reader.ReadLine()
		conn.WriteMessage(websocket.TextMessage, l)
	}
}
