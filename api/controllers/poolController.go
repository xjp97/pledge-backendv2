package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/api/models/response"
	"pledge-backendv2/api/services"
	"pledge-backendv2/api/validate"
	"pledge-backendv2/config"
	"regexp"
	"strings"
	"time"
)

type PoolController struct {
}

func (c *PoolController) PoolBaseInfo(ctx *gin.Context) {
	res := response.Gin{
		Res: ctx,
	}
	req := request.PoolBaseInfo{}
	var result []models.PoolBaseInfoRes
	// 校验参数
	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode = services.NewPool().PoolBaseInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

func (c *PoolController) TokenList(ctx *gin.Context) {
	req := request.TokenList{}
	result := response.TokenList{}

	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}
	errCode, data := services.NewTokenList().GetTokenList(&req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}
	var BaseUrl = c.GetBaseUrl()
	result.Name = "Pledge Token List"
	result.LogoURI = BaseUrl + "storage/img/Pledge-project-logo.png"
	result.Timestamp = time.Now()
	result.Version = response.Version{
		Major: 2,
		Minor: 16,
		Patch: 12,
	}
	for _, v := range data {
		result.Tokens = append(result.Tokens, response.Token{
			Name:     v.Symbol,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
			Address:  v.Token,
			ChainID:  v.ChainId,
			LogoURI:  v.Logo,
		})
	}
	ctx.JSON(200, result)
	return
}

func (c *PoolController) GetBaseUrl() string {

	domainName := config.Config.Env.DomainName
	domainNameSlice := strings.Split(domainName, "")
	pattern := "\\d+"
	isNumber, _ := regexp.MatchString(pattern, domainNameSlice[0])
	if isNumber {
		return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + ":" + config.Config.Env.Port + "/"
	}
	return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + "/"
}

func (c *PoolController) Search(ctx *gin.Context) {

	res := response.Gin{Res: ctx}
	req := request.Search{}
	result := response.Search{}

	// 参数校验
	errCode := validate.NewSearch().Search(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	err, tatol, pools := services.NewSearch().Search(&req)
	if err != statecode.CommonSuccess {
		res.Response(ctx, err, nil)
		return
	}
	result.Count = tatol
	result.Rows = pools
	res.Response(ctx, statecode.CommonSuccess, result)
	return

}

func (c *PoolController) DebtTokenList(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TokenList{}
	// 参数校验
	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	err, debtTokenList := services.NewTokenList().DebtTokenList(&req)
	if err != statecode.CommonSuccess {
		res.Response(ctx, err, nil)
		return
	}
	res.Response(ctx, statecode.CommonSuccess, debtTokenList)
	return
}

func (c *PoolController) PoolDataInfo(ctx *gin.Context) {
	res := response.Gin{
		Res: ctx,
	}
	req := request.PoolDataInfo{}
	var result []models.PoolDataInfoRes
	// 校验参数
	errCode := validate.NewPoolDataInfo().PoolDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode = services.NewPool().PoolDataInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
