package models

import "github.com/jinzhu/gorm"

type Quote struct {
	gorm.Model
	Symbol      string `json:"symbol"`
	CompanyName string `json:"company_name"`
}
