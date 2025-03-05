package routes

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/controllers"
	"pledge-backendv2/api/middlewares"
	"pledge-backendv2/config"
)

func InitRoute(e *gin.Engine) *gin.Engine {

	v2Group := e.Group("/api/v" + config.Config.Env.Version)

	poolController := controllers.PoolController{}
	// 查询借贷池基本信息
	v2Group.GET("/poolBaseInfo", poolController.PoolBaseInfo)
	v2Group.GET("poolDataInfo", poolController.PoolDataInfo)
	v2Group.GET("/token", poolController.TokenList)
	// 查询信息增加token校验
	v2Group.GET("/pool/debtTokenList", middlewares.CheckToken(), poolController.DebtTokenList)
	v2Group.GET("/pool/search", middlewares.CheckToken(), poolController.Search)

	userController := controllers.UserController{}
	v2Group.POST("/user/login", userController.Login)
	v2Group.POST("/user/logout", middlewares.CheckToken(), userController.Logout)

	priceController := controllers.PriceController{}
	v2Group.GET("/price", priceController.NewPrice)

	controller := controllers.MultiSignPoolController{}
	v2Group.POST("/pool/setMultiSign", middlewares.CheckToken(), controller.SetMultiSign)
	v2Group.POST("/pool/getMultiSign", middlewares.CheckToken(), controller.GetMultiSign)

	return e
}
