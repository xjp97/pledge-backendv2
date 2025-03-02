package services

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"pledge-backendv2/config"
	abifile "pledge-backendv2/contract/abi"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	"pledge-backendv2/schedule/models"
	"pledge-backendv2/utils"
	"strings"
)

type TokenSymbol struct{}

func NewTokenSymbol() *TokenSymbol {
	return &TokenSymbol{}
}

// 同步更新 合约符号
func (s *TokenSymbol) UpdateContractSymbol() {
	var tokens []models.TokenInfo
	db.Mysql.Table("token_info").Find(&tokens)
	for _, t := range tokens {
		if t.Token == "" {
			log.Logger.Sugar().Error("UpdateContractSymbol token empty", t.Symbol, t.ChainId)
			continue
		}
		err := errors.New("")
		symbol := ""
		if t.ChainId == config.Config.TestNet.ChainId {
			// 查询代币符号
			err, symbol = s.GetContractSymbolOnTestNet(t.Token, config.Config.TestNet.NetUrl)
		} else if t.ChainId == config.Config.MainNet.ChainId {
			// 主网络逻辑
		} else {
			log.Logger.Sugar().Error("UpdateContractSymbol chain_id err ", t.Symbol, t.ChainId)
			continue
		}
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractSymbol err ", t.Symbol, t.ChainId, err)
			continue
		}
		// 检查数据
		hasNewData, err := s.CheckSymbolData(t.Token, t.ChainId, symbol)
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractSymbol CheckSymbolData err ", t.Symbol, t.ChainId, err)
			continue
		}
		// 检查更新成功
		if hasNewData {
			// 更新表数据
			err = s.SaveSymbolData(t.Token, t.ChainId, symbol)
			if err != nil {
				log.Logger.Sugar().Error("UpdateContractSymbol SaveSymbolData err ", t.Symbol, t.ChainId, err)
				continue
			}
		}

	}

}

// 更新表符号数据
func (s *TokenSymbol) SaveSymbolData(token, chainId, symbol string) error {
	dateTimeFormat := utils.GetCurDateTimeFormat()
	// 更新表数据
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).Updates(map[string]interface{}{
		"symbol":     symbol,
		"updated_at": dateTimeFormat,
	}).Debug().Error
	if err != nil {
		return err
	}
	return nil
}

// 根据 token chain查询缓存, 数据库, 更新信息
func (s *TokenSymbol) CheckSymbolData(token, chainId, symbol string) (bool, error) {
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, err := db.RedisGet(redisKey)
	// 如果缓存中不存在
	if len(redisTokenInfoBytes) <= 0 {
		err = s.CheckTokenInfo(token, chainId)
		if err != nil {
			log.Logger.Error(err.Error())
		}
		db.RedisSet(redisKey, models.RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Symbol:  symbol,
		}, 0)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}
	} else {

		redisTokenInfo := models.RedisTokenInfo{}
		err = json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}
		// 如何符号相同, 直接返回
		if redisTokenInfo.Symbol == symbol {
			return false, nil
		}
		redisTokenInfo.Symbol = symbol

		err = db.RedisSet(redisKey, redisTokenInfo, 0)
		if err != nil {
			log.Logger.Error(err.Error())
			return true, err
		}
	}
	return true, nil
}

// 查询 token_info 表中是否有数据, 没有就新增
func (s *TokenSymbol) CheckTokenInfo(token, chainId string) error {
	tokenInfo := models.TokenInfo{}
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo = models.TokenInfo{}
			nowTime := utils.GetCurDateTimeFormat()
			tokenInfo.Token = token
			tokenInfo.ChainId = chainId
			tokenInfo.UpdatedAt = nowTime
			tokenInfo.CreatedAt = nowTime
			err = db.Mysql.Table("token_info").Create(&tokenInfo).Debug().Error
			if err != nil {
				return err
			}

		} else {
			return err
		}

	}
	return nil
}

// 通过合约地址查询合约符号
func (s *TokenSymbol) GetContractSymbolOnTestNet(token, network string) (error, string) {
	client, err := ethclient.Dial(network)
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	// 读取 erc20 abi文件
	abiStr, err := abifile.GetAbiByToken("erc20")
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	// 创建一个合约绑定对象
	contract := bind.NewBoundContract(common.HexToAddress(token), parsed, client, client, client)

	res := make([]interface{}, 0)
	// 调用合约 symbol 方法,获取代币符号
	err = contract.Call(nil, &res, "symbol")
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	return nil, res[0].(string)
}
