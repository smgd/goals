package models

import (
	"github.com/jinzhu/gorm"
)

type Area struct {
	gorm.Model
	Name        string
	Description string
	User        User `gorm:"foreignkey:UserID"`
	UserID      uint
}
