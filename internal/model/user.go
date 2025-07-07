package model

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Email             string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password          string    `gorm:"type:text;not null" json:"-"`
	IsVerified        bool      `gorm:"default:false" json:"is_verified"`
	VerificationToken string    `gorm:"type:varchar(255)" json:"-"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	UserType          string    `gorm:"type:enum('admin','standard');default:'standard'" json:"user_type"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

var DB *gorm.DB

func InitDB(user, pass, host, port, dbname string) error {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	DB = db
	return nil
}

func CreateUser(user *User) error {
	return DB.Create(user).Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	if user.ID == 0 {
		return errors.New("invalid user ID")
	}
	return DB.Save(user).Error
}

func GetAllUsers() ([]User, error) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
