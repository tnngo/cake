package cake

import (
	"fmt"

	"github.com/tnngo/cake/gmysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func loadGrom(configs []*mysqlConfig) {
	var (
		err error
	)

	var register *dbresolver.DBResolver
	for i, v := range configs {
		var dsn string
		if i == 0 {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", v.User, v.Password, v.IP, v.Port, v.Database)
			gmysql.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			if err != nil {
				zap.L().Error(err.Error(), zap.Reflect("config", v))
				return
			}
			continue
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", v.User, v.Password, v.IP, v.Port, v.Database)

		if register == nil {
			register = dbresolver.Register(dbresolver.Config{
				Sources: []gorm.Dialector{mysql.Open(dsn)},
			}, v.Database)
		} else {
			register = register.Register(dbresolver.Config{
				Sources: []gorm.Dialector{mysql.Open(dsn)},
			}, v.Database)
		}

		// if err := mysql.DB.Use(dbresolver.Register(dbresolver.Config{
		// 	Sources: []gorm.Dialector{gmysql.Open(dsn)},
		// }, v.Database)); err != nil {
		// 	zap.L().Error(err.Error(), zap.Reflect("config", v))
		// 	return
		// }

	}

	gmysql.DB.Use(register)
	sqlDB, _ := gmysql.DB.DB()
	sqlDB.SetMaxOpenConns(configs[0].MaxOpenConns)
	sqlDB.SetMaxIdleConns(configs[0].MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("数据库连接成功", zap.Reflect("configs", configs))
	}

}
