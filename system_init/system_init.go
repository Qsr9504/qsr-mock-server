package system_init

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"qsr-mock-server/common"
	"qsr-mock-server/models"
	"qsr-mock-server/service"
	"time"
)

func Init() {
	initConfig()
	//initMysql()
	service.Init() // 长链接初始化
}

func initConfig() {
	viper.SetConfigName(common.APPName)
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config app err:", err)
	}
	fmt.Println("config app init")
}

func initMysql() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL的日志
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色
		},
	)

	var err error
	dsnStr := viper.GetString("mysql.dsn")
	fmt.Println(dsnStr)
	common.DB, _ = gorm.Open(mysql.Open(dsnStr), &gorm.Config{Logger: newLogger}) // 并且配置自己的日志打印
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("MySQL Inited")
	// 数据表初试化
	if tableExist := common.DB.Migrator().HasTable(common.ProxyRulesModelTableName); !tableExist {
		// 表不存在就创建
		common.DB.AutoMigrate(&models.ProxyRulesModel{})
	}
	return common.DB
}
