package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mini-ebook/config"
	"mini-ebook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
