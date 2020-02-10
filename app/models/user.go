package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Username  string `gorm:"type:varchar(100);unique_index"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string
}
