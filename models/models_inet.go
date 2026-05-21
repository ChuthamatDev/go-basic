package models

import "gorm.io/gorm"

type Person struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	IsActive *bool  `json:"isactive" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=32"`
}

type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Color string `json:"color"`
}

type Register struct {
	Email        string `json:"email" validate:"required,email"`
	Username     string `json:"username" validate:"required,username_custom"`
	Password     string `json:"password" validate:"required,min=6,max=20"`
	LineID       string `json:"id_line,omitempty"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	BusinessType string `json:"business_type" validate:"required,business_type_custom"`
	Website      string `json:"website" validate:"required,website_custom"`
}

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

type ResultData struct {
	Data  []DogsRes `json:"data"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
}

type ResultDataV3 struct {
	Data        []DogsRes `json:"data"`
	Name        string    `json:"name"`
	Count       int       `json:"count"`
	SumRed      int       `json:"sum_red"`
	SumGreen    int       `json:"sum_green"`
	SumPink     int       `json:"sum_pink"`
	SumNoColor  int       `json:"sum_nocolor"`
}

func DogIDRangeScope(min, max int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("dog_id > ? AND dog_id < ?", min, max)
	}
}

func DeletedDogsScope(db *gorm.DB) *gorm.DB {
	return db.Unscoped().Where("deleted_at IS NOT NULL")
}

type Company struct {
	gorm.Model
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	TaxID       string `json:"tax_id"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}
