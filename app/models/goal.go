package models

import (
	"github.com/jinzhu/gorm"
)

type Goal struct {
	gorm.Model
	Name        string
	Description string
	Area        Area `gorm:"foreignkey:AreaID"`
	AreaID      uint
	User        User `gorm:"foreignkey:UserID"`
	UserID      uint
}
