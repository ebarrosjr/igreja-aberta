package database

import (
	"errors"
	"strings"
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

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

func MigrateUsers() {
	DB.AutoMigrate(&User{})
}

func CreateUser(congregationId int, name, email, password, status string) (*User, error) {
	if status == "" {
		status = UserStatusActive
	}

	user := &User{
		CongregationId: congregationId,
		Name:           strings.TrimSpace(name),
		Email:          strings.ToLower(strings.TrimSpace(email)),
		Password:       password,
		Status:         NormalizeUserStatus(status),
	}
	result := DB.Create(user)
	return user, result.Error
}

func ListUsers() ([]User, error) {
	var users []User
	result := DB.Order("id asc").Find(&users)
	return users, result.Error
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

func UserExistsByEmailExceptID(email string, id uint) (bool, error) {
	var count int64
	result := DB.Model(&User{}).Where("email = ? AND id <> ?", strings.ToLower(strings.TrimSpace(email)), id).Count(&count)
	return count > 0, result.Error
}

func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

func DeleteUser(user *User) error {
	return DB.Delete(user).Error
}

func NormalizeUserStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func IsValidUserStatus(status string) bool {
	switch NormalizeUserStatus(status) {
	case UserStatusActive, UserStatusInactive:
		return true
	default:
		return false
	}
}

func IsActiveStatus(status string) bool {
	return NormalizeUserStatus(status) == UserStatusActive
}
