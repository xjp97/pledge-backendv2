package models

import (
	"errors"
	"gorm.io/gorm"
	"pledge-backendv2/db"
	"pledge-backendv2/log"
	"pledge-backendv2/utils"
)

type PoolBase struct {
	Id                     int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	PoolId                 int    `json:"pool_id" gorm:"column:pool_id"`
	ChainId                string `json:"chain_id" gorm:"column:chain_id"`
	SettleTime             string `json:"settle_time" gorm:"column:settle_time"`
	EndTime                string `json:"end_time" gorm:"column:end_time"`
	InterestRate           string `json:"interest_rate" gorm:"column:interest_rate"`
	MaxSupply              string `json:"max_supply" gorm:"max_supply:"`
	LendSupply             string `json:"lend_supply" gorm:"column:lend_supply"`
	BorrowSupply           string `json:"borrow_supply" gorm:"column:borrow_supply"`
	MartgageRate           string `json:"martgage_rate" gorm:"column:martgage_rate"`
	LendToken              string `json:"lend_token" gorm:"column:lend_token"`
	LendTokenInfo          string `json:"lend_token_info" gorm:"column:lend_token_info"`
	BorrowToken            string `json:"borrow_token" gorm:"column:borrow_token"`
	BorrowTokenInfo        string `json:"borrow_token_info" gorm:"column:borrow_token_info"`
	State                  string `json:"state" gorm:"column:state"`
	SpCoin                 string `json:"sp_coin" gorm:"column:sp_coin"`
	JpCoin                 string `json:"jp_coin" gorm:"column:jp_coin"`
	LendTokenSymbol        string `json:"lend_token_symbol" gorm:"column:lend_token_symbol"`
	BorrowTokenSymbol      string `json:"borrow_token_symbol" gorm:"column:borrow_token_symbol"`
	AutoLiquidateThreshold string `json:"auto_liquidate_threshold" gorm:"column:auto_liquidate_threshold"`
	CreatedAt              string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt              string `json:"updated_at" gorm:"column:updated_at"`
}

type BorrowToken struct {
	BorrowFee  string `json:"borrowFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

type LendToken struct {
	LendFee    string `json:"lendFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

func NewPoolBase() *PoolBase {
	return &PoolBase{}
}

// 指定生成的表名称
func (p *PoolBase) TableName() string {
	return "poolbases"
}

func (p *PoolBase) SavePoolBase(chainId, poolId string, poolBase *PoolBase) error {
	// 获取当前时间
	nowDateTime := utils.GetCurDateTimeFormat()
	err, symbol := p.SaveTokenInfo(poolBase)
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	poolBase.BorrowTokenSymbol = symbol[0]
	poolBase.LendTokenSymbol = symbol[1]

	err = db.Mysql.Table("poobases").Where("chain_id=? and pool_id=?", chainId, poolId).First(&p).Debug().Error
	if err != nil {
		// 判断是否有记录,没有就创建
		if errors.Is(err, gorm.ErrRecordNotFound) {
			poolBase.CreatedAt = nowDateTime
			poolBase.UpdatedAt = nowDateTime
			err := db.Mysql.Table("poobases").Create(poolBase).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		} else {
			return errors.New("record select err" + err.Error())
		}
	}
	poolBase.UpdatedAt = nowDateTime

	err = db.Mysql.Table("poolbases").Where("chain_id=? and pool_id=?", chainId, poolId).Updates(poolBase).Debug().Error
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	return nil
}

// 根据 token 查询数据,返回对应的符号
func (p *PoolBase) SaveTokenInfo(base *PoolBase) (error, []string) {
	tokenInfo := TokenInfo{}
	tokenSymbol := []string{"", ""}
	dateTimeFormat := utils.GetCurDateTimeFormat()

	err := db.Mysql.Table("token_info").Where("chain_id=? and token=?", base.ChainId, base.BorrowToken).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo.Token = base.BorrowToken
			err := db.Mysql.Table("token_info").Create(&TokenInfo{
				ChainId:   base.ChainId,
				Token:     base.BorrowToken,
				CreatedAt: dateTimeFormat,
				UpdatedAt: dateTimeFormat,
			}).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err, tokenSymbol
			}
		} else {
			return errors.New("token_info record select err" + err.Error()), tokenSymbol
		}

	}
	// 设置符号
	tokenSymbol[0] = tokenInfo.Symbol
	tokenInfo = TokenInfo{}

	err = db.Mysql.Table("token_info").Where("chain_id=? and token=?", base.ChainId, base.LendToken).First(&tokenInfo).Debug().Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Mysql.Table("token_info").Create(&TokenInfo{
				Token:     base.LendToken,
				ChainId:   base.ChainId,
				CreatedAt: dateTimeFormat,
				UpdatedAt: dateTimeFormat,
			}).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err, tokenSymbol
			}
		} else {
			return errors.New("token_info record select err " + err.Error()), tokenSymbol
		}
	}
	tokenSymbol[1] = tokenInfo.Symbol
	return nil, tokenSymbol
}
