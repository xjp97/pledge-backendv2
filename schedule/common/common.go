package common

import (
	"github.com/joho/godotenv"
	"os"
	"pledge-backendv2/log"
)

var PlgrAdminPrivateKey string

func GetEnv() {

	var ok bool
	// 加载 env 配置文件
	if err := godotenv.Load(); err != nil {
		log.Logger.Warn("Error loading .env file")
	}
	// 读取配置文件key
	PlgrAdminPrivateKey, ok = os.LookupEnv("plgr_admin_private_key")
	if !ok {
		log.Logger.Error("environment variable is not set")
		panic("environment variable is not set")
	}
	//	fmt.Println("plgr_admin_private_key:", PlgrAdminPrivateKey)
}
