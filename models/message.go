package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 消息
type Message struct {
	gorm.Model
	UserId     int64  //发送者
	TargetId   int64  //接受者
	Type       int    //发送类型  1私聊  2群聊  3心跳
	Media      int    //消息类型  1文字 2表情包 3语音 4图片 /表情包
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

func (m Message) MarshalBinary() (data []byte, err error) {
	//TODO implement me
	return json.Marshal(m)
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}

// 映射关系
var clientMap = make(map[int64]*Node)

// 线程安全
var rwLocker sync.RWMutex

// 需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1.  获取参数 并 检验 token 等合法性
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isValida := true //checkToke()  待.........
	conn, _ := (&websocket.Upgrader{
		// 加入 token 校验逻辑
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)

	//2.获取conn
	currentNow := uint64(time.Now().Unix())
	// 创建连接节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}

	//3. 用户关系
	//4. userid 跟 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)
	//7.加入在线用户到缓存

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("data ========>", data)
		broadMsg(data)

	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 74, 176, 186),
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 向信道中发送数据
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc data: ", string(data))
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		var buf [1024]byte
		size, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecv: ", buf[0:size])

		// 调用后端数据处理逻辑
		dispatch(buf[0:size])
	}
}

// Init 初始化启动 udp 接受和发送方
func Init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("udp port goroutine inited!")
}

func dispatch(data []byte) {
	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	// 将 data 解码到 msg 中
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 私信
		fmt.Println("dispatch  data :", string(data))
		sendMsg(msg.TargetId, data)
	case 2: // 群发
		sendGroupMsg(msg.TargetId, data)
	}

}

func sendGroupMsg(id int64, data []byte) {

}

func sendMsg(tid int64, data []byte) {
	rwLocker.RLocker()
	node := clientMap[tid]
	rwLocker.RUnlock()
	// 处理数据
	msg := Message{}
	json.Unmarshal(data, &msg)

	targetId := tid
	userId := msg.UserId
	targetIdStr := strconv.Itoa(int(targetId))
	userIdStr := strconv.Itoa(int(userId))
	msg.CreateTime = uint64(time.Now().Unix())

	ctx := context.Background() // 创建初始会话，用于 Redis 操作
	result, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}
	if result != "" {
		fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(data))
		node.DataQueue <- data
	}
	var key string
	if targetId > userId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	r, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(r)) + 1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result() //jsonMsg
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}
