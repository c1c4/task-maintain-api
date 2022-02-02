package migration

import (
	"api/app/database"
	"api/app/models"
)

func getModels() []interface{} {
	return []interface{}{&models.User{}, &models.Task{}}
}

func AutoMigration() {
	database.Database.AutoMigrate(getModels()...)
}
