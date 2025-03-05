package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
)

type User struct{}

func NewUser() *User {
	return &User{}
}

func (u *User) Login(c *gin.Context, req *request.Login) int {
	// 将请求参数填充到结构体中
	err := c.ShouldBind(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "Name" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
			if e.Field() == "Password" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
		}
		return statecode.CommonErrServerErr
	}

	return statecode.CommonSuccess

}
