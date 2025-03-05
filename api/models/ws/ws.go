package ws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"pledge-backendv2/api/models/kucoin"
	"pledge-backendv2/config"
	"pledge-backendv2/log"
	"sync"
	"time"
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
var UserPingPongDurTime = config.Config.Env.WssTimeoutDuration // seconds

func (s *Server) ReadAndWrite() {
	errChan := make(chan error)
	Manager.Servers.Store(s.Id, s)

	defer func() {
		Manager.Servers.Delete(s.Id)
		_ = s.Socket.Close()
		close(s.Send)
	}()

	// 读取数据
	go func() {
		for {
			select {
			case msg, ok := <-s.Send:
				if !ok {
					errChan <- errors.New("write  message error")
					return
				}
				s.SendToClient(string(msg), SuccessCode)
			}
		}
	}()
	// 读取客户端心跳
	go func() {
		for {
			_, message, err := s.Socket.ReadMessage()
			if err != nil {
				log.Logger.Sugar().Error(s.Id+"ReadMessage err", err)
				errChan <- err
				return
			}
			if string(message) == "ping" || string(message) == `"ping"` || string(message) == "'ping'" {
				s.LastTime = time.Now().Unix()
				s.SendToClient("pong", PongCode)
			}
		}
	}()

	// 检查心跳
	for {
		select {
		case <-time.After(time.Second):
			if time.Now().Unix()-s.LastTime >= UserPingPongDurTime {
				s.SendToClient("heartbeat", ErrorCode)
				return
			}
		case err := <-errChan:
			log.Logger.Sugar().Error(s.Id+"ReadMessage err", err)
			return
		}
	}

}

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
