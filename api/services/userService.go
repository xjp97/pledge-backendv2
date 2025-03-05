package services

import (
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/api/models/response"
	"pledge-backendv2/config"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	"pledge-backendv2/utils"
)

type UserService struct{}

func NewUser() *UserService {
	return &UserService{}
}

// 登录校验用户, 增加缓存
func (s *UserService) Login(req *request.Login, result *response.Login) int {

	log.Logger.Sugar().Info("contractService", req)

	if req.Name == "admin" && req.Password == "password" {
		token, err := utils.CreateToken(req.Name)
		if err != nil {
			log.Logger.Sugar().Error(err)
			return statecode.CommonErrServerErr
		}
		result.TokenId = token
		// 设置缓存
		_ = db.RedisSet(req.Name, "login_ok", config.Config.Jwt.ExpireTime)
		return statecode.CommonSuccess
	}
	return statecode.NameOrPasswordErr

}
