package seed

import (
	"fmt"
	"interview1710/api/config"
	"interview1710/api/models"
)

func CreateTable() {
	if err := config.Database.Migrator().AutoMigrate(&models.Category{}, &models.SiteInfo{}); err != nil {
		fmt.Println("error while creating table", err)
	}
	fmt.Println("Table Created")
}
