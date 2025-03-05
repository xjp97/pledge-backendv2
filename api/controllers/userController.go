package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/api/models/response"
	"pledge-backendv2/api/services"
	"pledge-backendv2/api/validate"
	"pledge-backendv2/db"
)

type UserController struct {
}

func (u *UserController) Login(c *gin.Context) {
	res := response.Gin{
		c,
	}
	req := request.Login{}
	result := response.Login{}

	errCode := validate.NewUser().Login(c, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(c, errCode, nil)
		return
	}
	errCode = services.NewUser().Login(&req, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(c, errCode, nil)
		return
	}

	res.Response(c, statecode.CommonSuccess, &result)
	return
}

func (u *UserController) Logout(c *gin.Context) {
	res := response.Gin{Res: c}
	usernameIntf, _ := c.Get("username")
	_, _ = db.RedisDelete(usernameIntf.(string))
	res.Response(c, statecode.CommonSuccess, nil)
	return
}
