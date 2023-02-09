package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

const (
	Header = iota
	Body
	Uri
)

var (
	wg                sync.WaitGroup
	qsrWssServer      WssServer
	maxConn           = 160 // 最大可以存储多少个key
	maxChanLen        = 10
	deviceUpMesChan   chan *ChanMes
	deviceDownMesChan chan *ChanMes
	remoteUpMesChan   chan *ChanMes
	remoteDownMesChan chan *ChanMes
	remoteAddr        = "wss://t1-im-gateway.dewu.net/spider-service/v1/channels"
)

type ChanMes struct {
	MapKey  string
	MsgType int
	Msg     []byte
}

type MockPairClients struct {
	closed       bool
	RemoteClient *websocket.Conn // 消息发送方的conn
	DeviceClient *websocket.Conn // 消息接收方的conn
}

type WssServer struct {
	ClientPairMap  map[string]MockPairClients //将所有的ws对象都放进去,key 为，远程client + 终端client 的字符串拼接
	clientConnLock sync.RWMutex               // 加入Map读写锁
}

var upGrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() {
	// 将全局的Wss管理对象进行实例化
	qsrWssServer = WssServer{
		ClientPairMap: make(map[string]MockPairClients, maxConn),
	}
	deviceUpMesChan = make(chan *ChanMes, maxChanLen)
	deviceDownMesChan = make(chan *ChanMes, maxChanLen)
	remoteUpMesChan = make(chan *ChanMes, maxChanLen)
	remoteDownMesChan = make(chan *ChanMes, maxChanLen)

	go openMockWssListen() // 启动长链接转发协程
}

func GetKey(remoteClient, deviceClient *websocket.Conn) string {
	return remoteClient.LocalAddr().String() + deviceClient.RemoteAddr().String()
}

func (server *WssServer) getClient(key string) MockPairClients {
	var mockPairClients MockPairClients
	// 加锁，进行读取
	server.clientConnLock.Lock()
	// 取值
	mockPairClients = server.ClientPairMap[key]
	// 解锁
	server.clientConnLock.Unlock()
	return mockPairClients
}

func (server *WssServer) addClient(deviceClient, remoteClient *websocket.Conn) {
	// 加锁，进行读取
	server.clientConnLock.Lock()
	// 设置值
	server.ClientPairMap[GetKey(remoteClient, deviceClient)] = MockPairClients{
		RemoteClient: remoteClient,
		DeviceClient: deviceClient,
	}
	// 解锁
	server.clientConnLock.Unlock()
}

// removePairClient 删除成功返回true，不存在这个key返回false
func (server *WssServer) removePairClient(key string) bool {
	server.clientConnLock.Lock()
	if server.getClient(key).RemoteClient == nil {
		return false
	}
	delete(server.ClientPairMap, key)
	server.clientConnLock.Unlock()
	return true
}

// NewClientConn ************ service 对接处理的方法 ************
func NewClientConn(ctx *gin.Context) {
	//升级get请求为webSocket协议
	wsConn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	bs := []byte(string("welcome to qsr mock tcp server !!!"))
	fmt.Println(string(bs))
	wsConn.WriteMessage(websocket.TextMessage, bs)
	// 短连接升级长链接成功，需要立刻同步创建

	// 新建一个链接远程服务器
	dialer := websocket.Dialer{}
	remoteIMClient, _, err := dialer.Dial(remoteAddr, nil)
	if err != nil || remoteIMClient == nil {
		fmt.Println("创建远程客户端失败")
		wsConn.Close()
		return
	}
	fmt.Println("1. 添加了一个链接远程IM服务的 客户端：", remoteIMClient.LocalAddr().String())

	// 创建一个 MockPairClients 对象
	pairClients := MockPairClients{}
	pairClients.RemoteClient = remoteIMClient
	pairClients.DeviceClient = wsConn
	pairClients.closed = false // 新创建的是没有关闭的

	// 使用方法存入map
	qsrWssServer.addClient(wsConn, remoteIMClient)
	fmt.Println("2. 存入map中一对数据：", pairClients)
	go readRemoteIMClientMes(&pairClients)
	readDeviceIMClientMes(&pairClients)
	safeDelPairClientConn(&pairClients)
}

// 核心监听器
func openMockWssListen() {
	fmt.Println("openMockWssListen 启动整个wss核心监听")
	for {
		select {
		case v := <-deviceUpMesChan:
			fmt.Println("openMockWssListen 核心监听收到 设备客户端发来的：", string(v.Msg))
			// 终端发来的消息，给到上行消息处理器处理
			upMsgHandler(v)
		case v := <-remoteUpMesChan:
			// 直接使用远程客户端进行写给IM服务
			qsrWssServer.ClientPairMap[(*v).MapKey].RemoteClient.WriteMessage(v.MsgType, v.Msg)
		case v := <-remoteDownMesChan:
			// 链接远程的客户端 接收到了消息
			downMsgHandler(v)
		case v := <-deviceDownMesChan:
			// 使用终端客户端对象 写给终端
			fmt.Println("openMockWssListen 写给设备客户端消息：", string(v.Msg))
			qsrWssServer.ClientPairMap[(*v).MapKey].DeviceClient.WriteMessage(v.MsgType, v.Msg)
		}
	}
}

// 上行消息处理器
func upMsgHandler(upChanMsg *ChanMes) {
	fmt.Println("上行数据：\n", string(upChanMsg.Msg))
	// TODO:
	// 1. 将标记信息进行剥离
	// 2. 从数据库中根据标记信息获取指定的json串，获取成功，消息体进行替换；获取失败，直接上传消息
	// 3. 寻找与远程服务器连接的那个客户端对象，并将 确认后的 上行数据写出
	remoteUpMesChan <- upChanMsg
}

// 下行消息处理器
func downMsgHandler(downChanMsg *ChanMes) {
	// TODO:
	// 1. 从数据库中根据 下行数据拦截规则 获取指定的json串，获取成功，消息体进行替换；获取失败，直接上传消息
	fmt.Println("下行数据：\n", string(downChanMsg.Msg))
	// 2. 将 确认后的 下行数据写出到下行管道
	deviceDownMesChan <- downChanMsg
}

// ******** 读取 RemoteClient 的接收缓冲区，接收到了远程IM服务的消息返回 ************
func readRemoteIMClientMes(pairClients *MockPairClients) {
	// 启动读取两个client数据协程
	fmt.Println("启动读取 RemoteClient 数据协程")
	defer func() {
		fmt.Println("关闭 ---- readRemoteIMClientMes")
		pairClients.closed = true
		err := recover()
		if err != nil {
			fmt.Println("[ERROR]err=", err)
		}
	}()
	// 执行监听操作
	for {
		if pairClients.closed {
			fmt.Println("DeviceClient 已经关闭了,请求关闭该协程")
			return
		}
		remoteClientMessageType, remoteClientMessage, err := pairClients.RemoteClient.ReadMessage()
		fmt.Println("readRemoteIMClientMes 远程客户端收到了消息：", string(remoteClientMessage))
		if err != nil {
			fmt.Println("remoteIMClient Error during message reading:", err)
			return
		}
		if remoteClientMessage != nil {
			// 将消息移交给 下行管道
			remoteDownMesChan <- &ChanMes{
				MapKey:  GetKey(pairClients.RemoteClient, pairClients.DeviceClient),
				MsgType: remoteClientMessageType,
				Msg:     remoteClientMessage,
			}
		}
	}
}

// ******** 读取 DeviceClient 的接收缓冲区，接收到了远程IM服务的消息返回 ************
func readDeviceIMClientMes(pairClients *MockPairClients) {
	// 启动读取两个client数据协程
	fmt.Println("启动读取 DeviceClient 数据协程")
	defer func() {
		fmt.Println("关闭 ---- readDeviceIMClientMes")
		pairClients.closed = true
		err := recover()
		if err != nil {
			fmt.Println("[ERROR]err=", err)
		}
	}()
	// 执行监听操作
	for {
		if pairClients.closed {
			fmt.Println("RemoteClient 已经关闭了,请求关闭该协程")
			return
		}
		deviceClientMessageType, deviceClientMessage, err := pairClients.DeviceClient.ReadMessage() // 读取的时候会根据空格或\n 进行分割
		fmt.Println("设备客户端发来给mock服务 消息：", string(deviceClientMessage))
		if err != nil {
			fmt.Println("deviceClient Error during message reading:", err)
			return
		}
		if deviceClientMessage != nil {
			// 将消息塞进去上行管道
			deviceUpMesChan <- &ChanMes{
				MapKey:  GetKey(pairClients.RemoteClient, pairClients.DeviceClient),
				MsgType: deviceClientMessageType,
				Msg:     deviceClientMessage,
			}
		}
	}
}

// safeDelPairClientConn 安全的移除当前一对客户端对象
func safeDelPairClientConn(pairClients *MockPairClients) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("[ERROR]err=", err)
		}
	}() // 捕获异常，打印出来
	exist := qsrWssServer.removePairClient(GetKey(pairClients.RemoteClient, pairClients.DeviceClient))
	if !exist {
		fmt.Println("查询不存在，直接没了哦")
		// 如果删除的时候，不存在，就直接返回了
		return
	}
	pairClients.closed = true
	fmt.Println("安全移除两个客户端")
}
