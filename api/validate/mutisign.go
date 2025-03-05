package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
)

type MutiSign struct{}

func NewMutiSign() *MutiSign {
	return &MutiSign{}
}

func (v *MutiSign) SetMutiSign(c *gin.Context, req *request.SetMultiSign) int {
	err := c.ShouldBind(req)
	// 校验参数
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// chainId 不能为空
			if e.Field() == "SpName" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
		}
		return statecode.CommonErrServerErr
	}
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}
	return statecode.CommonSuccess
}

func (v *MutiSign) GetMutiSign(c *gin.Context, req *request.GetMultiSign) int {
	err := c.ShouldBind(req)
	// 校验参数
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// chainId 不能为空
			if e.Field() == "SpName" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
		}
		return statecode.CommonErrServerErr
	}
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}
	return statecode.CommonSuccess
}
