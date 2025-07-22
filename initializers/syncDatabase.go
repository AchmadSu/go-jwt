package initializers

import "example.com/m/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
