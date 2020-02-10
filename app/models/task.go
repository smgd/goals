package models

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
	Goal        Goal `gorm:"foreignkey:GoalID"`
	GoalID      uint
}
