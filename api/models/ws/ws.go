package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"pledge-backendv2/api/models/kucoin"
	"pledge-backendv2/log"
	"sync"
)

const SuccessCode = 0
const PongCode = 1
const ErrorCode = -1

type Server struct {
	sync.Mutex
	Id       string
	Socket   *websocket.Conn
	Send     chan []byte
	LastTime int64 // last send time
}

type ServerManager struct {
	Servers    sync.Map
	Broadcast  chan []byte
	Register   chan *Server
	Unregister chan *Server
}

type Message struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

var Manager = ServerManager{}

func (s *Server) SendToClient(data string, code int) {
	s.Lock()
	defer s.Unlock()
	dataBytes, err := json.Marshal(Message{
		Code: code,
		Data: data,
	})
	// 发送文本数据
	err = s.Socket.WriteMessage(websocket.TextMessage, dataBytes)
	if err != nil {
		log.Logger.Sugar().Error(s.Id+" SendToClient err ", err)
	}
}

// 启动 wesocket , 发送价格
func StartServer() {
	log.Logger.Info("websocket server start")
	for {
		select {
		case price, ok := <-kucoin.PlgrPriceChan:
			if ok {
				Manager.Servers.Range(func(k, v interface{}) bool {
					v.(*Server).SendToClient(price, SuccessCode)
					return true
				})
			}
		}
	}

}
