package models

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
