package models

import "errors"

// user model
type User struct {
	ID       uint   `gorm:"primary_key;auto_increment"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Admin    bool   `gorm:"not null,default:false"`
	Results1 []Format1
	Results2 []Format2
}

// validate that user fields are not empty
func (u *User) ValidateUser() error {
	if u.Name == "" {
		return errors.New("Name cannot be empty")
	}
	if u.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if u.Password == "" {
		return errors.New("Password cannot be empty")
	}
	return nil
}
