package models

import (
	"app/library/log"
	"app/library/task"
	"app/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"sync"
	"time"
)

var _db *gorm.DB
var _initOnce sync.Once

//init db connection
func initDB() {
	_initOnce.Do(func() {
		sqlUrl := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.SqlUsername,
			setting.SqlPassword,
			setting.SqlHost,
			setting.SqlPort,
			setting.SqlDatabase,
		)
		config := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "c_",
				SingularTable: true, //使用单数表名
			},
		}
		if setting.SqlDebug >= 1 {
			config.Logger = logger.Default.LogMode(logger.Info)
		} else {
			config.Logger = logger.Default.LogMode(logger.Silent)
		}

		log.StdOut("init", "db.start", sqlUrl)
		if err := task.Retry("db conn", 5, time.Second, func() error {
			conn, er := gorm.Open(mysql.Open(sqlUrl), config)
			if er != nil || conn == nil {
				return er
			}
			sqlDB, er := conn.DB()
			if er != nil {
				return er
			}
			sqlDB.SetMaxIdleConns(setting.SqlPoolMaxIdle)
			sqlDB.SetMaxOpenConns(setting.SqlPoolMaxOpen)
			_db = conn
			return nil
		}); err != nil {
			log.StdFatal("init", "db.err", err)
		}

		initDBPreHeating()
		log.StdOut("init", "db.ready")
	})
}

//init db root user
func initDBPreHeating() {
	//sql pre create
	if sqlBytes, err := setting.EmbedLocal().ReadFile("local/create.sql"); err == nil {
		var sqlItems = strings.Split(string(sqlBytes), ";")
		sqlItems = sqlItems[:len(sqlItems)-1]
		for _, sqlItem := range sqlItems {
			if err = _db.Exec(sqlItem).Error; err != nil {
				log.StdWarning("init", "db.table.init.err", err)
			}
		}
	}
	//init users
	if userRoot, err := CreateUser(UserRoot); err != nil {
		log.StdFatal("init", "db.user.root", err)
	} else {
		log.StdOut("init", "db.user.root.password:", userRoot.Password)
	}
	if userTest, err := CreateUser(UserTest); err != nil {
		log.StdFatal("init", "db.user.test", err)
	} else {
		log.StdOut("init", "db.user.test.password:", userTest.Password)
	}
	//init cfg
	createDefaultCfg()
}

//如果传了db连接，使用传入的db连接（用于事务开启场景）
func DB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) != 0 && tx[0] != nil {
		return tx[0]
	}
	return _db
}
