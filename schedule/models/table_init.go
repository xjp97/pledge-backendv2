package models

import "pledge-backendv2/db"

func InitTable() {
	db.Mysql.AutoMigrate(&PoolBase{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&RedisTokenInfo{})
	db.Mysql.AutoMigrate(&TokenInfo{})
}
