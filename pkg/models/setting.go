package models

import "gorm.io/gorm"

// Setting --
type Setting struct {
	*gorm.Model

	Name  string `json:"name" gorm:"primaryKey;autoIncrement:false"`
	Value string `json:"value"`
}
