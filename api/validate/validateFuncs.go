package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"pledge-backendv2/db"
	"regexp"
	"strings"
	"sync"
)

var oneofValsCache = map[string][]string{}
var oneofValsCacheRWLock = sync.RWMutex{}
var splitParamsRegexString = `'[^']*'|\S+`
var splitParamsRegex = regexp.MustCompile(splitParamsRegexString)

// 判断密码是否合法
func IsPassword(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^[a-zA-Z0-9!@#￥%^&*]{6,20}$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}

// 判断手机号是否合法
func IsPhoneNumber(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^1[0-9]{10}$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}

// IsEmail 检查手机号码字段是否合法
func IsEmail(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}

// CheckUserNicknameLength 检查用户昵称长度是否合法,限制1-20个字符即可，不同用户的昵称可以重复,第一个字符不可以是空格
func CheckUserNicknameLength(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		//if isOk, _ := regexp.MatchString(`^\w{1,20}$`, fl.Field().Interface().(string)); isOk {
		if isOk, _ := regexp.MatchString(`^.{1,20}$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}

// CheckUserAccount /*
func CheckUserAccount(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^[A-Za-z][A-Za-z0-9]{5,19}$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
		return false
	}
	return false
}

type dataStruct struct {
	DataCount int //这个结构体用来保存查询到的记录条数
}

// 字段唯一性约束
func OnlyOne(fl validator.FieldLevel) bool {
	vals := parseOnlyOneParam(fl.Param())
	if len(vals) <= 0 {
		panic("OnlyOne parameter err")
	}
	tableName := vals[0]
	fieldName := vals[1]

	var data dataStruct
	sqlStr := fmt.Sprintf("`%s`=?", fieldName)
	db.Mysql.Table(tableName).Select("COUNT(*)").Where(sqlStr, fl.Field().Interface()).Scan(&data.DataCount)
	if data.DataCount > 0 {
		return false
	}
	// 小于0说明没有重复记录
	return true
}

// 解析参数
func parseOnlyOneParam(s string) []string {
	// 加锁
	oneofValsCacheRWLock.RLock()
	vals, ok := oneofValsCache[s]
	oneofValsCacheRWLock.RUnlock()
	if !ok {
		oneofValsCacheRWLock.RLock()
		vals = splitParamsRegex.FindAllString(s, -1)
		for i := 0; i < len(vals); i++ {
			vals[i] = strings.Replace(vals[i], "'", "", -1)
		}
		oneofValsCache[s] = vals
		oneofValsCacheRWLock.Unlock()
	}
	return vals
}
