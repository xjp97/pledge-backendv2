package middlewares

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/response"
	"pledge-backendv2/config"
	"pledge-backendv2/db"
	"pledge-backendv2/utils"
)

func CheckToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		res := response.Gin{Res: c}
		authCode := c.Request.Header.Get("authCode")
		username, err := utils.ParseToken(authCode, config.Config.Jwt.SecretKey)
		if err != nil {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}
		if username != config.Config.DefaultAdmin.Username {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		resByteArr, err := db.RedisGet(username)
		if string(resByteArr) != `"login_ok"` {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
