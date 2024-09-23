package moudle

import (
	"context"
	"encoding/json"
	"fmt"
	"github/qq900306ss/SecondProject/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 訊息
type Message struct {
	gorm.Model
	UserId     int64  //發送者
	TargetId   int64  //接受者
	GroupId    int64  //群組
	Type       int    //發送類型  1私人  群組  3心跳
	Media      int    //訊息類型  1文字 2表情包 語音 4圖片
	Content    string //訊息內容
	CreateTime uint64 //創建時間
	ReadTime   uint64 //讀取時間
	Pic        string
	Url        string
	Desc       string
	Amount     int //沒用到
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	Addr      string // 地址
	DataQueue chan []byte
	GroupSets set.Interface //群組集合
}

var clientMap map[int64]*Node = make(map[int64]*Node, 0) // 映射關係

var rwLocker sync.RWMutex //讀寫鎖

// 發送者 接收者 類型 媒體 內容 圖片 影片 描述 數據
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.獲取參數並 檢驗token 等合法性
	// query.Get("token")

	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalida := true //token 驗證 //等等checkToke()
	conn, err := (&websocket.Upgrader{
		//token 驗證
		CheckOrigin: func(r *http.Request) bool { //防止跨領域請求
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	//2. 獲取conn

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//3. 用戶關係
	//4.userid 跟 node 綁定 並加鎖

	rwLocker.Lock()
	clientMap[userId] = node

	fmt.Println("調用了clientMap :", userId)

	rwLocker.Unlock()

	//5.完成發送邏輯

	//6.完成接受邏輯 你要看到自己消息所以要有
	go sendProc(node)

	go recvProc(node)

	SetUserOnlineInfo("online_"+Id, []byte(Id), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)

	sendMsg(userId, []byte("歡迎進入聊天室roger"))

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[放入]sendProc >>>> msg :", string(data))
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
		msg := Message{}
		err = json.Unmarshal(data, &msg) //可以解析Json 格式的Data 進而轉換成go 格式存入msg中
		if err != nil {
			fmt.Println(err)
		}
		dispatch(data)
		broadMsg(data)
		fmt.Println("[ws] recvProc <<<<< ", string(data))

	}
}

var udpsendchan chan []byte = make(chan []byte, 1024) //發送廣告的通道

func broadMsg(data []byte) { //跟接受消息廣播出來沒啥差別
	udpsendchan <- data
	fmt.Println("broadMsg  data :", string(data))
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine :")

}

// 完成udp 數據發送協成 只是做回應給網頁
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255), //這裡192.168.0.255 這樣的地址，這是一個廣播地址，會將數據發送到同一子網內的所有設備。  內部網絡：在內部網絡環境中，通常使用私有 IP 地址（例如 192.168.x.x、10.x.x.x、172.16.x.x 到 172.31.x.x 範圍）。
		Port: viper.GetInt("port.udp"),
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendchan:
			fmt.Println("udpSendProc  data :", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// 完成udp數據接收協程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, //net.IPv4zero 是一個特殊的 IP 地址 0.0.0.0。這個地址表示伺服器將在所有本地網絡接口上進行監聽。這樣，伺服器可以接收來自任何網絡接口的 UDP 數據包。
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte

		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc  data :", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 後端調度羅傑處理
func dispatch(data []byte) {
	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私人
		fmt.Println("dispatch  data :", string(data))
		sendMsg(msg.TargetId, data)
	case 2: //群發

		sendGroupMsg(msg.GroupId, data)

		fmt.Println("dispatch  data (group):", string(data))

	}
}

func sendMsg(userId int64, msg []byte) {

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg) //解析json

	ctx := context.Background() //不需要取消或截止時間的上下文：當你確定在某個操作過程中不需要取消操作或者設置截止時間時，可以使用 context.Background()。
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix()) //創建時間

	// r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()

	fmt.Println("測試主要地方:", string(msg), "userId:", userId)
	if ok {
		fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(msg))

		node.DataQueue <- msg

	} else {
		fmt.Println("ok怎麼了???", ok)
	}
	var key string
	if userId > jsonMsg.UserId && jsonMsg.GroupId == 0 {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result() //ZRevRange 倒數排序
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1                                     //紀錄+1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result() //jsonMsg //key有序集合不存在會自己建立出來

	fmt.Println("jsonMsg.GroupId : ", jsonMsg.GroupId)

	fmt.Println("測試一下key", key)
	fmt.Println("測試一下ress", ress)
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //備用 後續有機會拓展 紀錄完整msg 有機會的話
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}

func (msg Message) MarshalBinary() ([]byte, error) { //當 Message 結構體實現了 MarshalBinary 方法，它就符合 BinaryMarshaler 接口的要求，並且可以被用於需要二進制序列化的場景，比如將數據發送到網絡中，或者儲存到二進制文件
	return json.Marshal(msg)
}

func RedisMsg(userIdA int64, userIdB int64, groupId int64, start int64, end int64, isRev bool) []string { //單純用來回傳指定redis內容
	// rwLocker.RLock()
	//node, ok := clientMap[userIdA]
	// rwLocker.RUnlock()
	//jsonMsg := Message{}
	//json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	groupIdStr := strconv.Itoa(int(groupId))

	var key string
	if userIdA > userIdB && groupId == 0 {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	if groupId != 0 {

		key = "groupmsg_" + groupIdStr

	}
	fmt.Println("這裡的key :", key)
	//key = "msg_" + userIdStr + "_" + targetIdStr
	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result()  //根據分數倒數

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //沒有找到
	}
	// 發送推送消息
	/**
	// 後台通過消息 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	// contact.TagetId = comId
	contact.Type = 2
	community := Community{}
	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "群組不存在"
	}
	utils.DB.Where("owner_id=? and target_id=? and type =2", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "加過了"
	} else {
		contact.TagetId = community.ID
		contact.Desc = "群組"
		utils.DB.Create(&contact)
		return 0, "加入成功"
	}

}

func sendGroupMsg(userId int64, msg []byte) { //userid在這是groupid

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg) //解析json
	ctx := context.Background()   //不需要取消或截止時間的上下文：當你確定在某個操作過程中不需要取消操作或者設置截止時間時，可以使用 context.Background()。
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	groupIdStr := strconv.Itoa(int(jsonMsg.GroupId))
	jsonMsg.CreateTime = uint64(time.Now().Unix()) //創建時間
	// r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	number := SearchGroupWho(uint(userId))
	var key string
	for _, n := range number { //找群組成員有誰

		rwLocker.RLock()
		node, ok := clientMap[int64(n)]
		rwLocker.RUnlock()

		fmt.Println("測試主要地方:", string(msg), "userId:", n)
		if ok {
			if int64(n) == int64(jsonMsg.UserId) { //不發給自己
				continue
			}
			fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(msg))

			node.DataQueue <- msg

		} else {
			fmt.Println("ok怎麼了???", ok)
		}

	}

	if userId > jsonMsg.UserId && jsonMsg.GroupId == 0 {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	if jsonMsg.GroupId != 0 {

		key = "groupmsg_" + groupIdStr

	}

	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result() //ZRevRange 倒數排序
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1                                     //紀錄+1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result() //jsonMsg //key有序集合不存在會自己建立出來

	fmt.Println("jsonMsg.GroupId : ", jsonMsg.GroupId)

	fmt.Println("測試一下key", key)
	fmt.Println("測試一下ress", ress)
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //備用 後續有機會拓展 紀錄完整msg 有機會的話
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}
