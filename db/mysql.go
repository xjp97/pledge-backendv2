package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pledge-backendv2/config"
	"pledge-backendv2/log"
	"time"
)

func InitMysql() {

	mysqlConfig := config.Config.Mysql
	log.Logger.Info("Init Mysql")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.UserName,
		mysqlConfig.Password,
		mysqlConfig.Address,
		mysqlConfig.Port,
		mysqlConfig.DbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // 数据库路径
		DefaultStringSize:         256,   // string 类型字段长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度, mysql5.6之前版本不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并创建的方式
		DontSupportRenameColumn:   true,  // 用 change 重命名列 MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前mysql版本自动配置
	}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 关闭复数表
			},
			SkipDefaultTransaction: true,
		})

	if err != nil {
		log.Logger.Panic(fmt.Sprintf("mysql connention error ===> %+v", err))
	}
	// 创建回调函数
	_ = db.Callback().Create().After("gorm:after_create").Register("after_create", After)
	_ = db.Callback().Query().After("gorm:after_query").Register("after_query", After)
	_ = db.Callback().Delete().After("gorm:after_delete").Register("after_delete", After)
	_ = db.Callback().Update().After("gorm:after_update").Register("after_update", After)
	_ = db.Callback().Row().After("gorm:row").Register("after_row", After)
	_ = db.Callback().Raw().After("gorm:raw").Register("after_raw", After)

	sqlDb, err := db.DB()
	if err != nil {
		log.Logger.Panic(fmt.Sprintf("db.DB ERR ===> %+v", err))
	}
	// 空闲连接数
	sqlDb.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	// 最大连接数
	sqlDb.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(mysqlConfig.MaxLifeTime) * time.Second)

	Mysql = db
}

// 记录 sql
func After(db *gorm.DB) {
	db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	fmt.Println(db.Statement.Vars)
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	log.Logger.Info(sql)
}
