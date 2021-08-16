package utils

import (
	"errors"

	"gorm.io/gorm"
)

// StringSlicePosition --
func StringSlicePosition(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

// StringSliceContains --
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// CreateOrUpdate --
func CreateOrUpdate(db *gorm.DB, model interface{}, where interface{}, update interface{}) (interface{}, error) {
	var result interface{}
	err := db.Model(model).Where(where).First(result).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if err = db.Model(model).Create(update).Error; err != nil {
			return nil, err
		}
	}

	if err = db.Model(model).Where(where).Updates(update).Error; err != nil {
		return nil, err
	}

	return update, nil
}
