package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) Create() error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	err := DB().Create(&User{Username: u.Username, Password: string(hashedPassword)}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u User) IsExist() bool {
	var tmp User
	err := DB().Table("users").Where("username = ?", u.Username).First(&tmp).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func (u User) Validate() bool {
	var tmp User
	DB().Table("users").Where("username = ?", u.Username).First(&tmp)

	err := bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(u.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}
