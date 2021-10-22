package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"interview1710/api/config"
	"net/http"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SiteInfo struct {
	gorm.Model
	Domain     string         `json:"domain,omitempty"`
	CategoryID int            `json:"category_id,omitempty"`
	Tags       datatypes.JSON `json:"tags,omitempty"`
}

func AllDomain() ([]SiteInfo, error) {
	siteInfo := []SiteInfo{}

	err := config.Database.Model(&SiteInfo{}).Find(&siteInfo).Error
	return siteInfo, err
}

func DeleteDomain(id string) error {
	if config.Database.Debug().Delete(&SiteInfo{}, id).RowsAffected == 0 { //check if already deleted
		return errors.New("invalid domain_id " + id)
	}
	return nil
}

func UpdateDomain(r *http.Request) (SiteInfo, error) {
	siteInfo := SiteInfo{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&siteInfo); err != nil {
		return siteInfo, err
	}
	//make sure only updated_at field will be updated
	if config.Database.Model(SiteInfo{}).Debug().
		Where("site_infos.id = ?", siteInfo.ID).
		Updates(SiteInfo{Domain: siteInfo.Domain, CategoryID: siteInfo.CategoryID, Tags: siteInfo.Tags}).RowsAffected == 0 {
		return siteInfo, errors.New("invalid Domain.id change other")
	}
	//update elastic search with id doc
	UpdateIndex(siteInfo)
	return siteInfo, nil
}

func CreateDomain(r *http.Request) (SiteInfo, error) {
	siteInfo, err := validateDomain(r)
	if err != nil {
		return siteInfo, err
	}
	config.Database.Create(&siteInfo)
	return siteInfo, nil
}

func DomainCategory(catId string, limitNum int, offsNum int) ([]SiteInfo, error) {
	siteInfos := []SiteInfo{}
	fmt.Println(limitNum)
	if err := config.Database.Debug().Joins("join categories on categories.id = site_infos.category_id").
		Where("site_infos.category_id=?", catId).
		Limit(limitNum).
		Offset(offsNum).
		Find(&siteInfos).Error; err != nil {
		return siteInfos, err
	}
	return siteInfos, nil
}

func CreateDomainBasedCatName(inputSite SiteInfo) (SiteInfo, error) {
	site := SiteInfo{}
	if r, err := checkExistedDomain(inputSite.Domain); !r && err != nil {
		return site, err
	}
	config.Database.Create(&inputSite)
	return site, nil
}

func checkExistedDomain(domain string) (bool, error) {
	site := SiteInfo{}

	if config.Database.
		Where("domain = ?", domain).
		Limit(1).
		Find(&site).
		RowsAffected != 0 {
		return false, errors.New("domain " + domain + " has been created before ")
	}

	if config.Database.Error != nil {
		return false, errors.New("error checking siteinfore.domain ")
	}
	return true, nil
}

func validateDomain(r *http.Request) (SiteInfo, error) {
	siteInfo := SiteInfo{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&siteInfo); err != nil {
		return siteInfo, err
	}
	switch {
	case siteInfo.Domain == "":
		return siteInfo, errors.New("domain is empty, check it")
	case siteInfo.CategoryID == 0:
		return siteInfo, errors.New("categoryid is empty, check it")
	case siteInfo.Tags == nil:
		return siteInfo, errors.New("tag is empty, check it")
	}
	if r, err := checkExistedDomain(siteInfo.Domain); !r && err != nil {
		return siteInfo, err
	}
	return siteInfo, nil
}
