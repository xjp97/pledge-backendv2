package services

import (
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models"
	"pledge-backendv2/log"
)

type PoolService struct {
}

func NewPool() *PoolService {
	return &PoolService{}
}

func (s *PoolService) PoolBaseInfo(chainId int, result *[]models.PoolBaseInfoRes) int {
	err := models.NewPoolBases().PoolBaseInfo(chainId, result)
	if err != nil {
		log.Logger.Info(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess

}

func (s *PoolService) PoolDataInfo(chainId int, result *[]models.PoolDataInfoRes) int {
	err := models.NewPoolData().PoolDataInfo(chainId, result)
	if err != nil {
		log.Logger.Info(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}
