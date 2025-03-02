package main

import (
	"pledge-backendv2/db"
	"pledge-backendv2/schedule/models"
	"pledge-backendv2/schedule/task"
)

func main() {

	db.InitMysql()

	db.InitRedis()
	// 初始化表信息
	models.InitTable()

	task.Task()
}
