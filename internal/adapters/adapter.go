package adapters

import (
	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"gorm.io/gorm"
)

type CompanyAdapter struct {
	DB *gorm.DB
}

func NewCompanyAdapter(db *gorm.DB) *CompanyAdapter {
	return &CompanyAdapter{
		DB: db,
	}
}

func (company *CompanyAdapter) InsertCompanyCredentials(req entities.Credentials) (entities.Credentials, error) {

	query := "INSERT INTO credentials (company_id,company_username,name,aim,type_id) VALUES($1,$2,$3,$4,$5) RETURNING company_id,company_username,name,aim,type_id"
	var res entities.Credentials

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := company.DB.Raw(query, req.CompanyID,req.CompanyUsername,req.Name,req.Aim,req.TypeID).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}
	return res, nil
}

func (company *CompanyAdapter) InsertEmail(req []entities.CompanyEmail) ([]entities.CompanyEmail, error) {

	var res []entities.CompanyEmail

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := company.DB.Create(&req).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res,err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}
	return res, nil
}

func (company *CompanyAdapter) InsertPhone(req []entities.CompanyPhone) ([]entities.CompanyPhone, error) {

	var res []entities.CompanyPhone

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := company.DB.Create(&req).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res,err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}
	return res, nil
}

func (company *CompanyAdapter) InsertAddress(req entities.CompanyAddress) (entities.CompanyAddress, error) {

	var res entities.CompanyAddress

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := company.DB.Create(&req).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res,err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}
	return res, nil
}
