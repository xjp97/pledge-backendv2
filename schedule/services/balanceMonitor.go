package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"math/big"
	"pledge-backendv2/config"
	"pledge-backendv2/log"
	"pledge-backendv2/utils"
)

type BalanceMonitor struct {
}

func NewBalanceMonitor() *BalanceMonitor {
	return &BalanceMonitor{}
}

func (s *BalanceMonitor) Monitor() {
	// 查询代币余额
	balance, err := s.GetBalance(config.Config.TestNet.NetUrl, config.Config.TestNet.PledgePoolToken)

	thresholhPoolToken, ok := new(big.Int).SetString(config.Config.Threshold.PledgePoolTokenThresholdBnb, 10)
	// 如果余额小于 1 ,则发送邮件提醒
	if ok && (err == nil) && (balance.Cmp(thresholhPoolToken) <= 0) {
		emailBody, err := s.EmailBody(config.Config.TestNet.PledgePoolToken, "TBNB", balance.String(), thresholhPoolToken.String())
		if err != nil {
			log.Logger.Error(err.Error())
		} else {
			err := utils.SendEmail(emailBody, 2)
			if err != nil {
				log.Logger.Error(err.Error())
			}
		}
	}

}

// 封装邮件信息
func (s *BalanceMonitor) EmailBody(token, currency, balance, threshold string) ([]byte, error) {
	e18, err := decimal.NewFromString("1000000000000000000")

	if err != nil {
		return nil, err
	}
	balanceDeci, err := decimal.NewFromString(balance)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceDeci.Div(e18).String()
	thresholdDeci, err := decimal.NewFromString(threshold)
	if err != nil {
		return nil, err
	}
	thresholdStr := thresholdDeci.Div(e18).String()
	log.Logger.Sugar().Info("balance not enough", token, " ", currency, " ", balanceStr, " ", thresholdStr)
	body := fmt.Sprintf(`<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;The balance of <strong><span style="color: rgb(255, 0, 0);"> %s </span></strong> is <strong>%s %s </strong>. Please recharge it in time. The current minimum balance limit is %s %s 
</p>`, token, balanceStr, currency, thresholdStr, currency)
	return []byte(body), nil
}

// 查询 erc20 代币余额
func (s *BalanceMonitor) GetBalance(netUrl, token string) (*big.Int, error) {
	client, err := ethclient.Dial(netUrl)

	if err != nil {
		log.Logger.Error(err.Error())
		return big.NewInt(0), err
	}
	defer client.Close()
	// 根据 地址获取代币余额
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(token), nil)
	if err != nil {
		log.Logger.Error(err.Error())
		return big.NewInt(0), err
	}
	return balance, err
}
