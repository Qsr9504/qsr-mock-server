# 项目结构
```
.
├── docs
│   ├── docs.go
│   ├── README.md
│   ├── swagger.json
│   └── swagger.yaml
├── common  // 公共部分，包含常量封装
├── config  // 系统整体配置部分，包含数据库信息等 -> 待修改为配置中心
├── model   // 实体类，以及当前结构体直接对应的一些常用数据库操作
├── router  // 系统所有的路由地址信息 
├── service // 业务处理层
├── system_init // 系统初始化所需要的所有   
├── test    // 单元测试合集
├── utils    // 工具类封装
├── go.mod
├── go.sum
├── main.go
└── README.md
```


# 短连接mock配置方法：
1. 装配mock服务地址：实际远程ip地址前，加上mock服务的主机地址+ /mock（如：10.120.111.110:8080, 需要修改为 localhost:8787/mock/10.120.111.110:8080）
2. 添加拦截方法：通过页面添加一个拦截规则
满足拦截规则的就会返回指定返回值，不满足的就进行转发并原原本本返回

# 长链接mock配置方法：
1. 将客户端短连接进行mock，设置为mock服务主机的长链接地址
2. 添加mock规则（（包含关系、key、value）可以重复多个、预置的json或者直接截获不转发）
   - 模拟服务端推送：
     - 1）针对客户端发来的消息中，包含xxx=xxx的，就进行拦截，返回预置好的json串
     - 2）针对客户端发来的消息中

# 页面操作
------

## 一、背景
当前市面上有十分多的短连接mock平台，对于前端单测是极其方便的。
得益于后端架构绝大多数都是分布式，也有统一化的配置中心进行管理，域名会放在同一位置，因此微服务之间的HTTP类型接口的调用调试上，Mock服务都能发挥极大地作用。
常见的Mock平台有：
1. postman的本地mock
2. 公司内部的mooncake
3. Eolink
4. Apifox
.... 

都是关于 HTTP 或 HTTPS 请求进行mock，不依赖服务端的搭建
使用方式也比较固定，只需要编辑好固定的response JSON串，就可以快速启动一个指定路由的接口。

改进版：
一个路由地址只有一个返回串，不能够满足业务场景使用了。比如两个用户，使用同一个路由进行请求，但是请求头中的用户id或者其他字段是不同的，此时希望针对这两个不同用户有着不同的返回值。
因此，衍生了，请求头拦截规则、请求体拦截规则...  如果不满足拦截规则就进行放行，满足的话就返回指定的JSON信息。

长链接Mock：
长链接常见于IM和一些音视频相关的业务场景，市场上没有长链接mock平台。

## 二、业务痛点&应用场景
1. 业务场景链路深，长短链接结合，最终触发了一个定制化的长链接消息。
2. 前端性能的验证。往往需要大量的数据储备和触发动作，才能验证到一些隐藏比较深的前端场景
3. 长链接相关造数。
4. 长链接流量录制与回放。一些前端消息顺序导致的错误问题难以复现。将其录制之后，可以随时回放。

## 三、详细设计
### 1、长链接MockServer  实现
![image](https://user-images.githubusercontent.com/17244253/220855219-62f027a7-2ac0-40c6-a765-89e3df7588db.png)
1. 客户端链接上Mock Server，并发送消息
2. 设备服务器启动协程，专门服务这个「设备客户端」，将其接收到所有消息塞进去「设备客户端上行消息管道」
3. 「上行消息处理器」不断消费「设备客户端上行消息管道」的消息，并在此处做出上行消息拦截判断。
  - 如果此处拦截成功，意味着消息不再继续往后，整个消息到此消失。最上层的「IM Server」是不会接收到
  - 如果此处拦截失败，「上行消息处理器」会将消息继续塞进去「远程客户端上行消息管道」
4. 如果拦截失败/消息替换，消息流转到「远程客户端上行消息管道」
5. 「远程客户端」与「IM server」长链接相连，不断消费「远程客户端上行消息管道」消息
6. 然后将消息递交到「IM server」
7. 「远程客户端」拿到「IM server」的消息返回
8. 「远程客户端」将获取到的所有消息塞进去「远程客户端下行消息管道」
9. 「下行消息处理器」不断消费「远程客户端下行消息管道」，并且在此处做出下行消息拦截判断
  - 拦截成功，完成消息的替换或拦截
  - 拦截失败，消息继续下行
10. 拦截失败/消息替换，「下行消息处理器」将处理后的消息塞进「设备客户端下行消息管道」
11. 「设备服务端」持续消费「设备客户端下行消息管道」
12. 「设备服务端」将消息写出到「设备客户端」，完成整个长链接的消息链路

### 2、长链接拦截规则 实现
#### 1）动态规则配置
基于缓存实现的动态规则配置，即在当前会话中有效，服务重启或
> JSON 消息中包含 #*- （intercept rule） -*#
就是将规则拼接成一个JSON字符串，放入一个定制化的字符之间，如果包含这样的字符串，就动态的将规则添加进 redis 中，设置一个过期时间，就实现了指定时间段内规则生效。
> 可以实现的效果：在章鱼工作台中，输入 【#*- （订单卡片） -*#】，就可以让Mock服务直接模拟客户端推送过来一个「商品卡片」「订单卡片」消息。

- 消息处理器的业务流程图如下：
![image](https://user-images.githubusercontent.com/17244253/220855669-6cc065a1-1edc-4cfe-b744-e7b18052bce3.png)

#### 2）静态规则配置
设计原则：一个人使用Mock服务，尽可能不影响到其他人
理想的情况是这两个用户在使用Mock Server的时候，应该是要互不干涉的，并且再多一个用户C，不想使用Mock服务的时候，是需要真实访问到实际的目标ip地址的，拿到真实的返回数据的。
![image](https://user-images.githubusercontent.com/17244253/220855770-c6f89c55-09b5-486e-b97a-b771c8221a4c.png)
也就是所有人同时使用Mock服务，让想用mock的用，不想用的直接穿透mock进行放行。

### 4、短连接拦截规则  实现
**拦截规则实现**
我们知道所有的请求到达Mock Server之后，Mock Server都做了哪些事情呢？
![image](https://user-images.githubusercontent.com/17244253/220855929-5ce927dd-6765-4142-8524-ea243f48540d.png)

- 拦截实体
```go
type ProxyCondition struct {
   InterceptSite  int    // 拦截的目标位置，Header、Body、Uri 这几种方式
   IsContain      bool   // 包含、不包含
   InterceptKey   string // 拦截的 目标 key
   InterceptValue string // 拦截的 目标 value
}
type ProxyRulesModel struct {
   gorm.Model
   OwnerName          string           // 所属的用户
   ProxyMethod        string           // 需要代理的方法
   ProxyHost          string           // 需要代理的主机地址 ip + port
   ProxyUrl           string           // 需要代理的路由地址，不包含请求参数 和 http:// 或 https://
   ProxyConditionJson string           // 代理条件，json字符串，可以传入多个条件
   ProxyConditions    []ProxyCondition gorm:"-" // 代理的条件集合, gorm 忽略此字段
   CustomizedResp     string           // 定制的返回值
}
```
- 拦截核心逻辑
```go
// ProxyHandler
// @Summary    代理核心方法
// @Tags   代理模块
// @Router /mock/* [ANY]
func ProxyHandler(ctx *gin.Context) {
   remoteHost, remotePathWithOutParams := getHostAndPathFromCtx(ctx)
   // 从请求头中获取当前用户信息
   ownerName := ctx.Request.Header.Get("ownerName")
   fmt.Println("* mock服务 \n 解析出来的 remoteHost=", remoteHost, " remotePathWithOutParams=", remotePathWithOutParams, "\n mock服务 *\n")
   // 1. 根据 路由信息，从数据库中查询数据
   rulesModel := models.GetRulesByUserAndUrl(ownerName, remotePathWithOutParams)
   //fmt.Println("---checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson)---", checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson))
   // 2. 查询规则存在并且满足了拦截条件，才进行自定义mock数据返回
   if rulesModel.ProxyUrl != "" && checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson) {
      fmt.Println("--- 拦截成功 ---")
      jsonMap, err := utils.JsonStringToMap(rulesModel.CustomizedResp) // 将自定义的返回值 字符串 捞出来转换为map
      if err != nil {
         fmt.Println("JSON 字符串 转 Map 出错，检查数据库中 字符串是否是后边有人改动过")
         ctx.JSON(http.StatusInternalServerError, err)
         return
      }
      ctx.JSON(http.StatusOK, jsonMap) // 将自定义的返回值 设置为返回值
      return
   } else {
      // 拦截失败，直接转发，使用url中解析出来的远程ip地址转发
      var proxyUrl, _ = url.Parse(remoteHost)
      proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
      ctx.Request.URL.Path = remotePathWithOutParams // 修改上下文的路由地址
      ctx.Request.URL.Host = remoteHost              // 修改 上下文 中的 请求 主机地址，可以是host 也可以是 host:port
      proxy.ServeHTTP(ctx.Writer, ctx.Request)
      ctx.Abort()
   }
}
```
- 放行转发实现
```go
// 拦截失败，直接转发，使用url中解析出来的远程ip地址转发
var proxyUrl, _ = url.Parse(remoteHost)
proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
ctx.Request.URL.Path = remotePathWithOutParams // 修改上下文的路由地址
ctx.Request.URL.Host = remoteHost              // 修改 上下文 中的 请求 主机地址，可以是host 也可以是 host:port
proxy.ServeHTTP(ctx.Writer, ctx.Request)
ctx.Abort()
```
### 拦截与放行的视频演示：
视频演示的前置：设置了拦截规则，如果 请求体 中 不包含 qsr_key = 123 时，进行拦截，包含时则放行。
https://user-images.githubusercontent.com/17244253/220856252-9c3417dc-7350-4004-a740-094a053268bc.mp4

## 四、未来规划
1. 写个配套的前端页面
2. 把 RPC Mock 加进去（Dubbo、gRPC）
3. 长短连接进行强互动，即时间编排。如某一个接口触发后长链接推送指定的JSON消息
4. 搭配性能服务，对前后端性能进行场景压测。
如一个场景中使用ABC三个接口，此时，我们将D接口进行mock，使其不影响ABC接口的性能验证
1. 多返回值随机设置，目前mock server针对拦截到的请求，只允许预置一个返回消息。未来期望可以设置多个，并且进行随机（顺序）进行返回。

