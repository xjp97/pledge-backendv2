package models

import (
	"encoding/json"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/db"
)

type Pool struct {
	PoolID                 int      `json:"pool_id"`
	SettleTime             string   `json:"settleTime"`
	EndTime                string   `json:"endTime"`
	InterestRate           string   `json:"interestRate"`
	MaxSupply              string   `json:"maxSupply"`
	LendSupply             string   `json:"lendSupply"`
	BorrowSupply           string   `json:"borrowSupply"`
	MartgageRate           string   `json:"martgageRate"`
	LendToken              string   `json:"lendToken"`
	LendTokenSymbol        string   `json:"lend_token_symbol"`
	BorrowToken            string   `json:"borrowToken"`
	BorrowTokenSymbol      string   `json:"borrow_token_symbol"`
	State                  string   `json:"state"`
	SpCoin                 string   `json:"spCoin"`
	JpCoin                 string   `json:"jpCoin"`
	AutoLiquidateThreshold string   `json:"autoLiquidateThreshold"`
	Pooldata               PoolData `json:"pooldata"`
}

func NewPool() *Pool {
	return &Pool{}
}

func (p *Pool) Search(req *request.Search, whereCondition string) (error, int64, []Pool) {

	var tatol int64
	var pools []Pool          // 返回结构体
	var poolbases []PoolBases // 表实体

	// 查询总数
	db.Mysql.Table("poolbases").Where(whereCondition).Count(&tatol)
	// 查询数据
	err := db.Mysql.Table("poolbases").Where(whereCondition).Order("pool_id desc").Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&poolbases).Debug().Error

	if err != nil {
		return err, 0, pools
	}
	for _, v := range poolbases {
		var pooldata PoolData // 表实体
		err = db.Mysql.Table("pooldata").Where("chain_id", req.ChainID).Find(&pooldata).Debug().Error
		if err != nil {
			return err, 0, pools
		}
		borrowTokenInfo := BorrowTokenInfo{}
		_ = json.Unmarshal(([]byte)(v.BorrowToken), &borrowTokenInfo)
		lendTokenInfo := LendTokenInfo{}
		_ = json.Unmarshal(([]byte)(v.LendToken), &lendTokenInfo)

		pools = append(pools, Pool{
			PoolID:                 v.PoolID,
			SettleTime:             v.SettleTime,
			EndTime:                v.EndTime,
			InterestRate:           v.InterestRate,
			MaxSupply:              v.MaxSupply,
			LendSupply:             v.LendSupply,
			BorrowSupply:           v.BorrowSupply,
			MartgageRate:           v.MartgageRate,
			BorrowToken:            borrowTokenInfo.TokenName,
			LendToken:              lendTokenInfo.TokenName,
			State:                  v.State,
			SpCoin:                 v.SpCoin,
			JpCoin:                 v.JpCoin,
			AutoLiquidateThreshold: v.AutoLiquidateThreshold,
			Pooldata:               pooldata,
		})

	}
	return err, tatol, pools
}
