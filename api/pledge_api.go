package main

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/middlewares"
	"pledge-backendv2/api/models"
	"pledge-backendv2/api/models/kucoin"
	"pledge-backendv2/api/models/ws"
	"pledge-backendv2/api/routes"
	"pledge-backendv2/api/static"
	"pledge-backendv2/api/validate"
	"pledge-backendv2/config"
	"pledge-backendv2/db"
)

func main() {

	db.InitMysql()

	db.InitRedis()

	models.InitTable()

	validate.BindingValidator()

	go ws.StartServer()

	go kucoin.GetExchangePrice()

	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	staticPath := static.GetCurrentAbPathByCaller()
	app.Static("/storage/", staticPath)
	app.Use(middlewares.Cors()) // 跨域
	routes.InitRoute(app)
	_ = app.Run(":" + config.Config.Env.Port)

}
