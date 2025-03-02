package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"math/big"
	"pledge-backendv2/config"
	"pledge-backendv2/contract/bindings"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	serviceCommon "pledge-backendv2/schedule/common"
	"pledge-backendv2/schedule/models"
	"pledge-backendv2/utils"
	"time"
)

type TokenPrice struct {
}

func NewTokenPrice() *TokenPrice {
	return &TokenPrice{}
}

// 更新合约价格
func (s *TokenPrice) UpdateContractPrice() {
	var tokens []models.TokenInfo
	db.Mysql.Table("token_info").Find(&tokens)
	for _, t := range tokens {
		var err error
		var price int64 = 0
		if t.Token == "" {
			log.Logger.Sugar().Error("UpdateContractPrice token empty ", t.Symbol, t.ChainId)
			continue
		}
		if t.ChainId == config.Config.TestNet.ChainId {
			// 查询合约价格
			err, price = s.GetTestNetTokenPrice(t.Token)
		}
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractPrice err ", t.Symbol, t.ChainId, err)
			continue
		}
		// 检查合约价格和缓存中价格是否一致, 不一致更新
		hasNewData, err := s.CheckPriceData(t.Token, t.ChainId, utils.Int64ToString(price))
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractPrice err ", t.Symbol, t.ChainId, err)
			continue
		}
		if hasNewData {
			// 更新数据库价格
			err = s.SavePriceData(t.Token, t.ChainId, utils.Int64ToString(price))
			if err != nil {
				log.Logger.Sugar().Error("UpdateContractPrice SavePriceData err ", err)
				continue
			}
		}
	}

}

func (s *TokenPrice) SavePlgrPriceTestNet() {
	price := 22222
	client, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	// 获取合约实例
	bscPledgeOracleTestnetToken, err := bindings.NewBscPledgeOracleTestnetToken(common.HexToAddress(config.Config.TestNet.BscPledgeOracleToken), client)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	privateKey, err := crypto.HexToECDSA(serviceCommon.PlgrAdminPrivateKey)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	// 创建交易签名
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(utils.StringToInt64(config.Config.TestNet.ChainId)))

	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	transactOpts := bind.TransactOpts{
		From:      auth.From,
		Nonce:     nil,
		Signer:    auth.Signer, // Method to use for signing the transaction (mandatory)
		Value:     big.NewInt(0),
		GasPrice:  nil,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  0,
		Context:   ctx,
		NoSend:    false, // Do all transact steps but do not send the transaction
	}
	_, err = bscPledgeOracleTestnetToken.SetPrice(&transactOpts, common.HexToAddress(config.Config.TestNet.PlgrAddress), big.NewInt(int64(price)))
	log.Logger.Sugar().Info("SavePlgrPrice", err)
	// 获取价格
	a, d := s.GetTestNetTokenPrice(config.Config.TestNet.PlgrAddress)
	fmt.Println(a, d, 5555)

}

// 更新对应 token 价格数据
func (s *TokenPrice) SavePriceData(token, chainId, price string) error {

	nowDateTime := utils.GetCurDateTimeFormat()
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=? ", token, chainId).Updates(map[string]interface{}{
		"price":      price,
		"updated_at": nowDateTime,
	}).Debug().Error
	if err != nil {
		log.Logger.Sugar().Error("UpdateContractPrice SavePriceData err ", err)
		return err
	}

	return nil
}

func (s *TokenPrice) CheckPriceData(token, chainId, price string) (bool, error) {
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, err := db.RedisGet(redisKey)
	if len(redisTokenInfoBytes) <= 0 {
		// 检查token是否存在,不存在新增
		err = s.CheckTokenInfo(token, chainId)
		if err != nil {
			log.Logger.Sugar().Error("CheckPriceData err ", token, chainId, err)
		}
		err = db.RedisSet(redisKey, models.RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Price:   price,
		}, 0)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}

	} else {
		redisTokenInfo := models.RedisTokenInfo{}
		err = json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			log.Logger.Sugar().Error(err.Error())
			return false, err
		}
		if redisTokenInfo.Price == price {
			return true, nil
		}
		redisTokenInfo.Price = price
		err = db.RedisSet(redisKey, redisTokenInfo, 0)
		if err != nil {
			log.Logger.Sugar().Error(err.Error())
			return true, err
		}

	}
	return true, nil
}

// 查询是否有数据, 如果没有就新增
func (s *TokenPrice) CheckTokenInfo(token, chainId string) error {

	tokenInfo := models.TokenInfo{}
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo := models.TokenInfo{}
			dateTimeFormat := utils.GetCurDateTimeFormat()
			tokenInfo.Token = token
			tokenInfo.ChainId = chainId
			tokenInfo.CreatedAt = dateTimeFormat
			tokenInfo.UpdatedAt = dateTimeFormat
			err := db.Mysql.Table("token_info").Create(&tokenInfo).Debug().Error
			if err != nil {
				log.Logger.Sugar().Error("CreateTokenInfo err ", tokenInfo.ChainId, err)
				return err
			} else {
				return err
			}
		}
	}
	return nil
}

func (s *TokenPrice) GetTestNetTokenPrice(token string) (error, int64) {
	client, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return err, 0
	}
	// 获取合约
	bscPledgeOracleTestnetToken, err := bindings.NewBscPledgeOracleTestnetToken(common.HexToAddress(config.Config.TestNet.BscPledgeOracleToken), client)
	if nil != err {
		log.Logger.Error(err.Error())
		return err, 0
	}
	price, err := bscPledgeOracleTestnetToken.GetPrice(nil, common.HexToAddress(token))
	if nil != err {
		log.Logger.Error(err.Error())
		return err, 0
	}
	return nil, price.Int64()

}
