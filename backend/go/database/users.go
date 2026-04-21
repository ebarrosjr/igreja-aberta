package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CongregationId int        `json:"congregation_id" gorm:"not null;index"`
	Name           string     `json:"name" gorm:"type:varchar(255);not null"`
	Email          string     `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Password       string     `json:"password" gorm:"type:varchar(255);not null"`
	Status         string     `json:"status" gorm:"type:varchar(50);not null"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func MigrateUsers() {
	DB.AutoMigrate(&User{})
}

func CreateUser(congregationId int, name, email, password string) (*User, error) {
	user := &User{
		CongregationId: congregationId,
		Name:           name,
		Email:          email,
		Password:       password,
		Status:         "active",
	}
	result := DB.Create(user)
	return user, result.Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	result := DB.First(&user, id)
	return &user, result.Error
}

func UserExistsByEmail(email string) (bool, error) {
	var count int64
	result := DB.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}
