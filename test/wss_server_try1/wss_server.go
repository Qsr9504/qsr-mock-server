package wss

import (
	"fmt"
	"net"
)

/*
为了接收多个客户端传来的消息，架起来的服务端
1. 移动端链接上mock服务，
*/

var (
	normalClient      = 1 // 与Mock服务，即当前服务进行连接的客户端
	forImServerClient = 2 // 与IM中台服务进行对接的客户端

	maxConn        = 160 // 最大可以存储多少个key
	maxBuffLen     = 10240 / 2
	clientConnMap  = make(map[string]*net.TCPConn, maxConn)
	connKeyPairMap = make(map[string]string, maxConn) // 将一对key保存两次保存进入map
)

// 监听所有客户端请求
func listenClient(serverAddr, remoteAddr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", serverAddr)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for { // 循环接收
		clientConn, _ := tcpListener.AcceptTCP()                               // 监听请求连接
		clientConnMap[clientConn.RemoteAddr().String()] = clientConn           // 将连接添加到 connMap,使用RemoteAddr().String()作为key，才是客户端的来源地址
		go addMsgReceiver(clientConn)                                          // 开启一个协程专门监听这个客户端的所有来信
		go connectRemoteIMServer(remoteAddr, clientConn.RemoteAddr().String()) // 开启一个协程专门链接远程服务
		fmt.Println("用户 : ", clientConn.RemoteAddr().String(), " 已连接.")
	}
}

// 增加一个普通客户端的消息监听器，mock服务接收 普通客户端消息
func addMsgReceiver(newConnect *net.TCPConn) {
	for {
		byteMsg := make([]byte, maxBuffLen)
		len, err := newConnect.Read(byteMsg) // 从newConnect中读取信息到缓存中
		if err != nil {
			newConnect.Close()
		}
		fmt.Println(string(byteMsg[:len]))
		msgHandler(byteMsg[:len], newConnect.RemoteAddr().String())
	}
}

// 所有的消息处理器,不管是什么类型的客户端都会 通过这个处理器进行处理
func msgHandler(byteMsg []byte, key string) {
	for k := range connKeyPairMap {
		if k == key { // 只转发给目标客户端
			clientConnMap[key].Write(byteMsg)
		}
	}
}

//--------- 连接远程IM服务的客户端相关操作 -----------

// 建立连接
func connectRemoteIMServer(remoteIMAddr, clientKey string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", remoteIMAddr) // 使用tcp
	con, err := net.DialTCP("tcp", nil, tcpAddr)          // 拨号：主动向server建立连接
	if err != nil {
		fmt.Println("连接服务器失败")
	}
	// 将连接远程IM服务的 客户端 放入 map集合
	clientConnMap[con.LocalAddr().String()] = con
	fmt.Println("添加了一个链接远程IM服务的 客户端：", con.LocalAddr().String())
	go addMsgReceiver(con)
}

// 启动长链接代理服务
func MockIMServerRun(serverAddr, remoteAddr string) {
	go listenClient(serverAddr, remoteAddr)
}

// 消息接收器
//func msgReceiver(con *net.TCPConn, clientKey string) {
//	fmt.Println("msgReceiver", utils.GetGid())
//	buff := make([]byte, maxBuffLen)
//	for {
//		mes, _ := con.Read(buff) // 从建立连接的缓冲区读消息
//		fmt.Println(string(buff[:mes]))
//	}
//}

// 消息发送器
//func msgSender(con *net.TCPConn, clientKey string) {
//	fmt.Println("msgSender", utils.GetGid())
//	for {
//		bMsg, _, _ := reader.ReadLine()
//		bMsg = []byte(clientKey + " : " + string(bMsg))
//		con.Write(bMsg) // 发消息
//	}
//}
