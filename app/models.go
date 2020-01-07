package app

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

type Area struct {
	gorm.Model
	Name        string
	Description string
	User        User `gorm:"foreignkey:UserID"`
	UserID      uint
}

type Goal struct {
	gorm.Model
	Name        string
	Description string
	Area        Area `gorm:"foreignkey:AreaID"`
	AreaID      uint
}

type Task struct {
	gorm.Model
	Name        string
	Description string
	Goal        Goal `gorm:"foreignkey:GoalID"`
	GoalID      uint
}
