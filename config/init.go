package config

import (
	"github.com/BurntSushi/toml"
	"path"
	"path/filepath"
	"runtime"
)

// 读取配置文件配置
func init() {
	currentAbPathByCaller := getCurrentAbPathByCaller()

	abs, err := filepath.Abs(currentAbPathByCaller + "/configV21.toml")
	if err != nil {
		panic("read config file error:" + err.Error())
		return
	}
	if _, err := toml.DecodeFile(abs, &Config); err != nil {
		panic("read config file error:" + err.Error())
		return
	}
}

func getCurrentAbPathByCaller() string {

	var abPath string
	_, fileName, _, ok := runtime.Caller(0)

	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
