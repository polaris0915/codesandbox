package model

import (
	"github.com/polaris/codesandbox/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDB *gorm.DB
)

func InitAllDB() {
	initMysqlDB()
}

func initMysqlDB() {
	// 处理mysql数据库连接
	// 连接本地数据库，数据库的信息应该写入配置文件
	dsn := "root:" + settings.MysqlConfig.MysqlDBPassword + "@tcp(127.0.0.1:3306)/" + settings.MysqlConfig.MysqlDBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("mysql数据库初始化失败, 请检查: " + err.Error())
	}
	MysqlDB = DB
}
