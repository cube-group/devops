package models

import (
	"app/library/crypt/md5"
	"app/library/log"
	"app/library/uuid"
	"app/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
)

const UserRoot = "root"

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
		conn, err := gorm.Open(mysql.Open(sqlUrl), config)
		if err != nil || conn == nil {
			log.StdFatal("init", "db.err", err)
		}
		if sqlDB, err := conn.DB(); err == nil {
			sqlDB.SetMaxIdleConns(setting.SqlPoolMaxIdle)
			sqlDB.SetMaxOpenConns(setting.SqlPoolMaxOpen)
		}
		_db = conn

		initDBPreHeating()
		log.StdOut("init", "db.ready")
	})
}

//init db root user
func initDBPreHeating() {
	//sql pre
	//if sqlBytes, err := ioutil.ReadFile("local/db.sql"); err == nil {
	//	var sqlItems = strings.Split(string(sqlBytes), ";")
	//	sqlItems = sqlItems[:len(sqlItems)-1]
	//	for _, sqlItem := range sqlItems {
	//		if err = _db.Exec(sqlItem).Error; err != nil {
	//			log.StdWarning("init", "db.table.init.err", err)
	//		}
	//	}
	//}

	var rootRandPwd string
	var user = new(User).Get(UserRoot)
	if user == nil {
		rootRandPwd = md5.MD5(uuid.GetUUID("devops.root"))
		var newUser = &User{
			Username: UserRoot,
			RealName: UserRoot,
			Adm:      1,
		}
		if err := newUser.PwdAddUser(rootRandPwd); err != nil {
			log.StdFatal("init", "db.user.root.err", err.Error())
		}
		user = newUser
	}
	log.StdOut("init", "db.user.root.password:", user.Password)
}

//如果传了db连接，使用传入的db连接（用于事务开启场景）
func DB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) != 0 {
		return tx[0]
	}
	return _db
}
