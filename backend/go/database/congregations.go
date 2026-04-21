package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Congregation struct {
	Id           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"type:varchar(255);not null"`
	Description  string     `json:"description" gorm:"type:varchar(1000)"`
	Phone        string     `json:"phone" gorm:"type:varchar(30)"`
	Email        string     `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	ZipCode      string     `json:"zip_code" gorm:"type:varchar(20)"`
	Address      string     `json:"address" gorm:"type:varchar(255)"`
	Number       string     `json:"number" gorm:"type:varchar(20)"`
	Complement   string     `json:"complement" gorm:"type:varchar(255)"`
	Neighborhood string     `json:"neighborhood" gorm:"type:varchar(255)"`
	City         string     `json:"city" gorm:"type:varchar(255)"`
	State        string     `json:"state" gorm:"type:varchar(100)"`
	Status       string     `json:"status" gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}

func Migrate() {
	DB.AutoMigrate(&Congregation{})
}

func CreateCongregation(name, description, phone, email, zipCode, address, number, complement, neighborhood, city, state string) (*Congregation, error) {
	congregation := &Congregation{
		Name:         name,
		Description:  description,
		Phone:        phone,
		Email:        email,
		ZipCode:      zipCode,
		Address:      address,
		Number:       number,
		Complement:   complement,
		Neighborhood: neighborhood,
		City:         city,
		State:        state,
		Status:       "active",
	}
	result := DB.Create(congregation)
	return congregation, result.Error
}

func GetAllCongregations() ([]Congregation, error) {
	var congregations []Congregation
	result := DB.Find(&congregations)
	return congregations, result.Error
}

func GetCongregationByID(id uint) (*Congregation, error) {
	var congregation Congregation
	result := DB.First(&congregation, id)
	return &congregation, result.Error
}

func GetCongregationByEmail(email string) (*Congregation, error) {
	var congregation Congregation
	result := DB.Where("email = ?", email).First(&congregation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &congregation, nil
}
