package utils

import (
	"fmt"
	"strconv"
	"time"
)

// 格式化时间字符串
func FormatDate(timestamp string) string {
	var intTimeStamp, err2 = strconv.ParseInt(timestamp, 10, 64)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	return time.Unix(intTimeStamp, 0).Format("2006-01-02 15:04:05")
}

// 判断该权限是否存在
func CheckPower(power string, powerList []string, userId int) bool {
	// 超级管理员可访问
	if userId == 1 {
		return true
	}
	for _, v := range powerList {
		if power == v {
			return true
		}
	}
	return false
}

// 判断该模块是否存在
func CheckModule(module int, moduleList []string, userId int) bool {
	// 超级管理员可访问
	if userId == 1 {
		return true
	}
	for _, v := range moduleList {
		if strconv.Itoa(module) == v {
			return true
		}
	}
	return false
}
