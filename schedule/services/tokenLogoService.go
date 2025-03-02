package services

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"pledge-backendv2/config"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	"pledge-backendv2/schedule/models"
	"pledge-backendv2/utils"
)

type TokenLogo struct{}

func NewTokenLogo() *TokenLogo {
	return &TokenLogo{}
}

func (s *TokenLogo) UpdateTokenLogo() {

	res, err := utils.HttpGet(config.Config.Token.LogoUrl, map[string]string{})
	if err != nil {
		log.Logger.Sugar().Info("UpdateTokenLogo HttpGet error:", err)
	} else {
		tokenLogoRemote := models.TokenLogoRemote{}
		err = json.Unmarshal(res, &tokenLogoRemote)
		if err != nil {
			log.Logger.Sugar().Info("UpdateTokenLogo json.Unmarshal error:", err)
			return
		}
		for _, t := range tokenLogoRemote.Tokens {
			hasNewData, err := s.CheckLogoData(t.Address, utils.IntToString(t.ChainID), t.LogoURI, t.Symbol)
			if err != nil {
				log.Logger.Sugar().Info("UpdateTokenLogo CheckLogoData error:", err)
				continue
			}
			if hasNewData {
				err := s.SaveTokenData(t.Address, utils.IntToString(t.ChainID), t.LogoURI, t.Symbol, t.Decimals)
				if err != nil {
					log.Logger.Sugar().Info("UpdateTokenLogo SaveTokenData error:", err)
					continue
				}
			}
		}
	}

}

// 更新表数据
func (s *TokenLogo) SaveTokenData(token, chainId, logoUrl, symbol string, decimals int) error {
	dateTimeFormat := utils.GetCurDateTimeFormat()

	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).Updates(map[string]interface{}{
		"symbol":     symbol,
		"decimals":   decimals,
		"logo":       logoUrl,
		"updated_at": dateTimeFormat}).Debug().Error

	if err != nil {
		log.Logger.Sugar().Info("SaveTokenData error:", err)
		return err
	}
	return nil
}

// 检查更新缓存数据库
func (s *TokenLogo) CheckLogoData(token, chainId, logoUrl, symbol string) (bool, error) {
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, err := db.RedisGet(redisKey)

	if len(redisTokenInfoBytes) <= 0 {
		// 检查  网络下token是否存在, 如何不存在创建
		err = s.CheckTokenInfo(token, chainId)
		if err != nil {
			log.Logger.Sugar().Info("CheckTokenInfo error:", err)
		}
		err = db.RedisSet(redisKey, models.RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Logo:    logoUrl,
			Symbol:  symbol,
		}, 0)
		if err != nil {
			log.Logger.Sugar().Info("CheckTokenInfo redis error:", err)
			return false, err
		}
	} else {
		redisTokenInfo := models.RedisTokenInfo{}
		err = json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			log.Logger.Sugar().Info("CheckTokenInfo error:", err)
			return false, err
		}
		if redisTokenInfo.Logo == token {
			return false, nil
		}
		redisTokenInfo.Logo = logoUrl
		redisTokenInfo.Symbol = symbol

		err = db.RedisSet(redisKey, redisTokenInfo, 0)
		if err != nil {
			log.Logger.Sugar().Info("CheckTokenInfo  RedisSet error:", err)
			return true, err
		}

	}
	return true, nil
}

func (s *TokenLogo) CheckTokenInfo(token, chainId string) error {
	tokenInfo := models.TokenInfo{}
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo = models.TokenInfo{}
			nowDateTime := utils.GetCurDateTimeFormat()
			tokenInfo.Token = token
			tokenInfo.ChainId = chainId
			tokenInfo.UpdatedAt = nowDateTime
			tokenInfo.CreatedAt = nowDateTime
			err = db.Mysql.Table("token_info").Create(&tokenInfo).Debug().Error
			if err != nil {
				log.Logger.Sugar().Info("CheckTokenInfo error:", err)
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
