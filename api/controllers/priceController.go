package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"pledge-backendv2/api/models/ws"
	"pledge-backendv2/log"
	"pledge-backendv2/utils"
	"strings"
	"time"
)

type PriceController struct{}

// 升级连接为 ws, 监听客户端消息
func (c *PriceController) NewPrice(ctx *gin.Context) {

	// 监听运行时异常
	defer func() {
		recoverRes := recover()
		if recoverRes != nil {
			log.Logger.Sugar().Error("new price recover err:", recoverRes)
		}
	}()

	// 连接升级, 创建ws连接
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger.Sugar().Error(err.Error())
		return
	}
	// randomId 设置客户端唯一标识
	randomId := ""
	remoteIP := ctx.RemoteIP()
	if remoteIP != "" {
		randomId = strings.Replace(remoteIP, ".", "_", -1) + "_" + utils.GetRandomString(23)
	} else {
		randomId = utils.GetRandomString(32)
	}

	server := &ws.Server{
		Id:       randomId,
		Socket:   conn,
		Send:     make(chan []byte, 800),
		LastTime: time.Now().Unix(),
	}

	go server.ReadAndWrite()

}
