package models

import "pledge-backendv2/db"

// 初始化表
func InitTable() {
	db.Mysql.AutoMigrate(&MultiSign{})
	db.Mysql.AutoMigrate(&TokenInfo{})
	db.Mysql.AutoMigrate(&TokenList{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&PoolBases{})
}
