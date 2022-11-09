package starter

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/global"
)

// # 初始化db链接
func initDb() {
	global.Db = gormMysql()
}

func gormMysql() *gorm.DB {
	m := config.Conf.Mysql
	if m == nil || m.Dbname == "" {
		global.Log.Panic("未找到数据库配置信息")
		return nil
	}
	global.Log.Infof("连接mysql [%s]", m.Host)
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), //DSN date source name
		DefaultStringSize:         191,     //string
		DisableDatetimePrecision:  true,    //禁用 datetime 精度, Mysql 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	ormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表命名前缀
			SingularTable: true,
		}, Logger: logger.Default.LogMode(logger.Info), // 改为logger.Info即可显示sql语句. #原来这个值是Error类型
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), ormConfig); err != nil {
		global.Log.Panicf("连接mysql失败！[%s]", err.Error())
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
