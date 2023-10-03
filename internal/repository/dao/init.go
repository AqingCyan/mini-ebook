package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	// 严格说，不是好的实践
	return db.AutoMigrate(&User{})
}
