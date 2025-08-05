package helpers

import (
	"errors"

	"example.com/m/initializers"
	"gorm.io/gorm"
)

func PreloadRelationByID(modelPtr any, id uint, relations []string) error {
	db := initializers.DB

	for _, rel := range relations {
		db = db.Preload(rel)
	}

	if err := db.First(modelPtr, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	return nil
}
