package models

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserId uint `gorm:"not_null;index"`
	Title string `gorm:"not_null"`
}