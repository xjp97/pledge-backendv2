package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/api/models/response"
	"pledge-backendv2/api/services"
	"pledge-backendv2/api/validate"
	"pledge-backendv2/log"
)

type MultiSignPoolController struct {
}

func (c *MultiSignPoolController) SetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.SetMultiSign{}
	log.Logger.Sugar().Info("SetMultiSign req ", req)

	errCode := validate.NewMutiSign().SetMutiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMutiSign().SetMultiSign(&req)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, nil)
	return
}

func (c *MultiSignPoolController) GetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.GetMultiSign{}
	result := response.MultiSign{}
	log.Logger.Sugar().Info("GetMultiSign req ", nil)

	errCode := validate.NewMutiSign().GetMutiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMutiSign().GetMultiSign(&result, req.ChainId)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
