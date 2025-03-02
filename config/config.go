package config

var Config Conf

type Conf struct {
	Mysql     MysqlConfig
	Redis     RedisConfig
	TestNet   TestNetConfig
	MainNet   MainNetConfig
	Token     TokenConfig
	Threshold ThresholdConfig
	Email     EmailConfig
}

type MysqlConfig struct {
	Address      string `toml:"address"`
	Port         string `toml:"port"`
	DbName       string `toml:"db_name"`
	UserName     string `toml:"user_name"`
	Password     string `toml:"password"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxLifeTime  int    `toml:"max_life_time"`
}

type RedisConfig struct {
	Address     string `toml:"address"`
	Port        string `toml:"port"`
	Db          int    `toml:"db"`
	UserName    string `toml:"user_name"`
	Password    string `toml:"password"`
	MaxIdle     int    `toml:"max_idle"`
	MaxActive   int    `toml:"max_active"`
	IdleTimeout int    `toml:"idle_timeout"`
}

type TestNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

type MainNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

type TokenConfig struct {
	LogoUrl string `toml:"logo_url"`
}

type ThresholdConfig struct {
	PledgePoolTokenThresholdBnb string `toml:"pledge_pool_token_threshold_bnb"`
}

type EmailConfig struct {
	Username string   `toml:"username"`
	Pwd      string   `toml:"pwd"`
	Host     string   `toml:"host"`
	Port     string   `toml:"port"`
	From     string   `toml:"from"`
	Subject  string   `toml:"subject"`
	To       []string `toml:"to"`
	Cc       []string `toml:"cc"`
}
