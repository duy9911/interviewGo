package models

import (
	"encoding/json"
	"errors"
	"interview1710/api/config"
	"net/http"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name      string     `json:"name,omitempty"`
	SiteInfos []SiteInfo `gorm:"foreignKey:CategoryID" json:"siteinfors,omitempty"`
}

func AllCategory() ([]Category, error) {
	categories := []Category{}

	err := config.Database.Model(&Category{}).Find(&categories).Error
	return categories, err
}

func DeleteCategory(id string) error {
	if config.Database.Debug().Delete(&Category{}, id).RowsAffected == 0 {
		return errors.New("invalid category_id " + id)
	}
	return nil
}

func CreateCategory(r *http.Request) (Category, error) {
	category, err := validateCategory(r)
	if err != nil {
		return category, err
	}

	//make sure only updated_at field will be updated
	if err := config.Database.Omit("ID", "Created_at").Create(&category).Error; err != nil {
		return category, errors.New("error create category " + category.Name)
	}
	return category, nil
}

func UpdateCategory(r *http.Request) (Category, error) {
	category, err := validateCategory(r)
	if err != nil {
		return category, err
	}

	//make sure only updated_at field will be updated
	if config.Database.Model(Category{}).Debug().
		Where("categories.id = ?", category.ID).
		Updates(Category{Name: category.Name}).RowsAffected == 0 {
		return category, errors.New("invalid id change other")
	}
	return category, nil
}

func CreateCategoryBasedName(name string) (Category, error) {
	category := Category{}
	//check exits or not
	if config.Database.
		Where("categories.name =?", name).
		Limit(1).
		Find(&category).
		RowsAffected != 0 {
		return category, errors.New("categories_name " + category.Name + " has been created before ")
	}

	if config.Database.Error != nil {
		return category, errors.New("error checking categories.name ")
	}

	input := Category{
		Name:      name,
		SiteInfos: nil,
	}
	//create categories.name
	if err := config.Database.Select("Name").Create(&input).Error; err != nil {
		return category, errors.New("create categories.name errors")
	}
	return category, nil
}

func OneCategoryName(name string) (Category, error) {
	category := Category{}
	if err := config.Database.Where("name = ?", name).Find(&category).Error; err != nil {
		return category, errors.New("rror checking name category")
	}
	return category, nil
}

func validateCategory(r *http.Request) (Category, error) {
	category := Category{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&category)

	if err != nil {
		return category, err
	}
	if category.Name == "" {
		return category, errors.New("category name is empyt")
	}

	if config.Database.
		Where("categories.name =?", category.Name).
		Limit(1).
		Find(&category).
		RowsAffected != 0 {
		return category, errors.New("categories_name " + category.Name + " has been created before ")
	}

	if config.Database.Error != nil {
		return category, errors.New("error checking categories.name ")
	}
	return category, nil
}
