package db

import (
	"auth-service/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMariaDB(user, pass, host, port, dbname string) error {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	model.DB = db // ✅ IMPORTANT: Make this assignment so model package can use the same DB

	if !db.Migrator().HasTable(&model.User{}) {
		if err := db.Migrator().CreateTable(&model.User{}); err != nil {
			return err
		}
	}

	return nil
}
