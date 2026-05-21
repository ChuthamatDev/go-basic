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
}

type ResultData struct {
	Data  []DogsRes `json:"data"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
}

type ResultDataV2 struct {
	Data     []DogsRes `json:"data"`
	NameDog  string    `json:"name_dog"`
	CountSum int       `json:"count_sum"`
}

type Company struct {
	gorm.Model
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	TaxID       string `json:"tax_id"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}
