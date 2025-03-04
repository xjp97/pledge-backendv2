package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
)

type Search struct{}

func NewSearch() *Search {
	return &Search{}
}

func (v *Search) Search(c *gin.Context, req *request.Search) int {
	err := c.ShouldBind(req)
	// 校验参数
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// chainId 不能为空
			if e.Field() == "ChainId" && e.Tag() == "required" {
				return statecode.ChainIdEmpty
			}
		}
		return statecode.CommonErrServerErr
	}
	if req.ChainID != 97 && req.ChainID != 56 {
		return statecode.ChainIdErr
	}
	return statecode.CommonSuccess
}
