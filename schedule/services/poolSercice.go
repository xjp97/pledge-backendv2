package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"pledge-backendv2/config"
	"pledge-backendv2/contract/bindings"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	serviceCommon "pledge-backendv2/schedule/common"
	"pledge-backendv2/schedule/models"
	"pledge-backendv2/utils"
	"strings"
	"time"
)

type poolService struct{}

// 创建一个 service 对象
func NewPool() *poolService {
	return &poolService{}
}

func (s *poolService) UpdateAllPoolInfo() {
	s.UpdatePoolInfo(config.Config.TestNet.PledgePoolToken, config.Config.TestNet.NetUrl, config.Config.TestNet.ChainId)
}

// 更新池信息
func (s *poolService) UpdatePoolInfo(contractAddress, network, chainId string) {

	log.Logger.Sugar().Info("UpdatePoolInfo", "contractAddress", contractAddress, "network", network, "chainId", chainId)
	client, err := ethclient.Dial(network)
	if err != nil {
		log.Logger.Sugar().Error("UpdatePoolInfo", "err", err)
		return
	}
	// 加载合约
	pledgePoolToken, err := bindings.NewPledgePoolToken(common.HexToAddress(contractAddress), client)
	// abi 生成的只读 go 绑定 PledgePoolTokenCaller
	// 获取合约费率
	borrowFee, err := pledgePoolToken.PledgePoolTokenCaller.BorrowFee(nil)
	fmt.Println(borrowFee)
	lendFee, err := pledgePoolToken.PledgePoolTokenCaller.LendFee(nil)
	fmt.Println(lendFee)

	privateKeyEcdsa, err := crypto.HexToECDSA(serviceCommon.PlgrAdminPrivateKey)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyEcdsa, big.NewInt(utils.StringToInt64(config.Config.MainNet.ChainId)))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactOpts := bind.TransactOpts{
		From:      common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF"),
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

	settleTime := big.NewInt(1709443200) // 示例Unix时间戳
	endTime := big.NewInt(1712035200)    // 示例Unix时间戳
	interestRate := uint64(500)          // 代表5%的年利率，假设是以千分之几的形式表示
	maxSupply := big.NewInt(1000)        // 示例最大供应量
	mortgageRate := big.NewInt(50)       // 示例抵押率，如50%
	lendTokenAddr := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	borrowTokenAddr := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	spTokenAddr := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	jpTokenAddr := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	autoLiquidateThreshold := big.NewInt(80) // 示例自动清算阈值

	createPoolInfo, err := pledgePoolToken.CreatePoolInfo(
		&transactOpts,
		settleTime,
		endTime,
		interestRate,
		maxSupply, mortgageRate, lendTokenAddr,
		borrowTokenAddr, spTokenAddr, jpTokenAddr, autoLiquidateThreshold)
	fmt.Println(createPoolInfo)

	poolLength, err := pledgePoolToken.PledgePoolTokenCaller.PoolLength(nil)
	if err != nil {
		log.Logger.Sugar().Error("UpdatePoolInfo", "err", err)
	}
	for i := 0; i <= int(poolLength.Int64())-1; i++ {
		log.Logger.Sugar().Info("UpdatePoolInfo ", i)
		poolId := utils.IntToString(i + 1)
		baseInfo, err := pledgePoolToken.PledgePoolTokenCaller.PoolBaseInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Error("UpdatePoolInfo", "err", err)
			continue
		}
		_, borrowToken := models.NewTokenInfo().GetTokenInfo(baseInfo.BorrowToken.String(), chainId)
		_, lendToken := models.NewTokenInfo().GetTokenInfo(baseInfo.LendToken.String(), chainId)

		lendTokenJson, _ := json.Marshal(models.LendToken{
			LendFee:    lendFee.String(),
			TokenLogo:  lendToken.Logo,
			TokenName:  lendToken.Symbol,
			TokenPrice: lendToken.Price,
		})

		borrowTokenJson, _ := json.Marshal(models.LendToken{
			LendFee:    borrowFee.String(),
			TokenLogo:  borrowToken.Logo,
			TokenName:  borrowToken.Symbol,
			TokenPrice: borrowToken.Price,
		})

		poolbase := models.PoolBase{
			SettleTime:             baseInfo.SettleTime.String(),
			PoolId:                 utils.StringToInt(poolId),
			ChainId:                chainId,
			EndTime:                baseInfo.EndTime.String(),
			InterestRate:           baseInfo.InterestRate.String(),
			MaxSupply:              baseInfo.MaxSupply.String(),
			LendSupply:             baseInfo.LendSupply.String(),
			BorrowSupply:           baseInfo.BorrowSupply.String(),
			MartgageRate:           baseInfo.MartgageRate.String(),
			LendToken:              baseInfo.LendToken.String(),
			LendTokenSymbol:        lendToken.Symbol,
			LendTokenInfo:          string(lendTokenJson),
			BorrowToken:            baseInfo.BorrowToken.String(),
			BorrowTokenSymbol:      borrowToken.Symbol,
			BorrowTokenInfo:        string(borrowTokenJson),
			State:                  utils.IntToString(int(baseInfo.State)),
			SpCoin:                 baseInfo.SpCoin.String(),
			JpCoin:                 baseInfo.JpCoin.String(),
			AutoLiquidateThreshold: baseInfo.AutoLiquidateThreshold.String(),
		}
		hasInfoData, byteBaseInfoStr, baseInfoMd5Str := s.GetPoolMd5(&poolbase, "base_info:pool_"+chainId+"_"+poolId)
		// 判断是否是新数据, 如果是
		if !hasInfoData || (baseInfoMd5Str != byteBaseInfoStr) {
			// 新增数据
			err := models.NewPoolBase().SavePoolBase(chainId, poolId, &poolbase)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolBase", "err", "chainId", "poolId", err, chainId, poolId)
			}
			// 添加借贷池缓存
			_ = db.RedisSet("base_info:pool_"+chainId+"_"+poolId, baseInfoMd5Str, 60*30)
		}

		// 查询借贷池结算数据
		dataInfo, err := pledgePoolToken.PledgePoolTokenCaller.PoolDataInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Info("UpdatePoolInfo PoolBaseInfo err", poolId, err)
			continue
		}
		hasPoolData, byteDataInfoStr, dataInfoMd5Str := s.GetPoolMd5(&poolbase, "base_data:pool_"+chainId+"_"+poolId)
		if !hasPoolData || (dataInfoMd5Str != byteDataInfoStr) { // have new data
			poolData := models.PoolData{
				PoolId:                 poolId,
				ChainId:                chainId,
				FinishAmountBorrow:     dataInfo.FinishAmountBorrow.String(),
				FinishAmountLend:       dataInfo.FinishAmountLend.String(),
				LiquidationAmounBorrow: dataInfo.LiquidationAmounBorrow.String(),
				LiquidationAmounLend:   dataInfo.LiquidationAmounLend.String(),
				SettleAmountBorrow:     dataInfo.SettleAmountBorrow.String(),
				SettleAmountLend:       dataInfo.SettleAmountLend.String(),
			}
			err := models.NewPoolData().SavePoolData(chainId, poolId, &poolData)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolData err ", chainId, poolId)
			}
			_ = db.RedisSet("data_info:pool_"+chainId+"_"+poolId, dataInfoMd5Str, 60*30)
		}

	}

}

// 将 数据转 md5 存入缓存
func (s *poolService) GetPoolMd5(baseInfo *models.PoolBase, key string) (bool, string, string) {
	// 转换成 json格式数组
	baseInfoBytes, _ := json.Marshal(baseInfo)
	// 将字节数组转 md5 加密
	baseInfoMd5Str := utils.Md5(string(baseInfoBytes))
	resInfoBytes, _ := db.RedisGet(key)
	if len(resInfoBytes) > 0 {
		return true, strings.Trim(string(resInfoBytes), `"`), baseInfoMd5Str
	}
	return false, strings.Trim(string(resInfoBytes), `"`), baseInfoMd5Str
}
