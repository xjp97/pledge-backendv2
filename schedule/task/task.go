package task

import (
	"github.com/jasonlvhit/gocron"
	"pledge-backendv2/db"
	"pledge-backendv2/schedule/common"
	"pledge-backendv2/schedule/services"
	"time"
)

func Task() {
	// 读取私钥
	common.GetEnv()
	// 清空当前db
	err := db.RedisFlushDB()
	if err != nil {
		panic("clear redis error" + err.Error())
	}
	// 初始化任务
	// 同步合约借贷池数据到mysql
	services.NewPool().UpdateAllPoolInfo()
	// 更新 token_info 表价格数据
	services.NewTokenPrice().UpdateContractPrice()
	// 同步代币符号
	services.NewTokenSymbol().UpdateContractSymbol()
	//同步 logo
	services.NewTokenLogo().UpdateTokenLogo()
	// 监控用户代币余额,小于 1的时候发送邮件提醒
	services.NewBalanceMonitor().Monitor()
	services.NewTokenPrice().SavePlgrPriceTestNet()

	s := gocron.NewScheduler()
	s.ChangeLoc(time.UTC)

	_ = s.Every(2).Minutes().From(gocron.NextTick()).Do(services.NewPool().UpdateAllPoolInfo)
	_ = s.Every(1).Minute().From(gocron.NextTick()).Do(services.NewTokenPrice().UpdateContractPrice)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenSymbol().UpdateContractSymbol)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenLogo().UpdateTokenLogo)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewBalanceMonitor().Monitor)
	//_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPrice)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPriceTestNet)
	<-s.Start()
}
