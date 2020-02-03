package models

import (
	"github.com/jinzhu/gorm"
)

type Area struct {
	gorm.Model
	Name        string
	Weight      int
	Description string
	Icon        string
	IsFavourite bool `gorm:"default:false"`
	User        User `gorm:"foreignkey:UserID"`
	UserID      uint
}
