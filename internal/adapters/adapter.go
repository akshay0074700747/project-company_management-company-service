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

	query := "INSERT INTO credentials (company_id,company_username,name,type_id) VALUES($1,$2,$3,$4) RETURNING company_id,company_username,name,type_id"
	var res entities.Credentials

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := company.DB.Raw(query, req.CompanyID, req.CompanyUsername, req.Name, req.TypeID).Scan(&res).Error; err != nil {
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
		return res, err
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
		return res, err
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
		return res, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}
	return res, nil
}

func (comp *CompanyAdapter) IsCompanyUsernameExists(companyUsername string) (bool, error) {

	query := "SELECT * FROM credentials WHERE company_username = $1"

	res := comp.DB.Exec(query, companyUsername)
	if res.Error != nil {
		return true, res.Error
	}

	if res.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (comp *CompanyAdapter) AttachCompanyRoleAndPermissions(req entities.CompanyRoles) error {

	query := "INSERT INTO company_roles (company_id,role_id,permission_id) VALUES($1,$2,$3)"

	if err := comp.DB.Exec(query, req.CompanyID, req.RoleID, req.PermissionID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) AddMember(req entities.CompanyMembers) error {

	query := "INSERT INTO company_members (company_id,role_id,member_id) VALUES($1,$2,$3)"

	if err := comp.DB.Exec(query, req.CompanyID, req.RoleID, req.MemberID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) IsMemberExists(id string) (bool, error) {

	query := "SELECT * FROM company_members WHERE member_id = $1"

	res := comp.DB.Exec(query, id)
	if res.Error != nil {
		return true, res.Error
	}

	if res.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (comp *CompanyAdapter) IsRoleIDExists(roleID uint) (bool, error) {

	query := "SELECT * FROM company_roles WHERE id = $1"

	res := comp.DB.Exec(query, roleID)
	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (comp *CompanyAdapter) GetRoleWithPermissionIDs(companyID string) ([]entities.CompanyRoles, error) {

	query := "SELECT * FROM company_roles WHERE company_id = $1"
	var res []entities.CompanyRoles

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetPermissions() ([]entities.Permissions, error) {

	query := "SELECT * FROM permissions"
	var res []entities.Permissions

	if err := comp.DB.Raw(query).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetCompanyTypes() ([]entities.CompanyTypes, error) {

	query := "SELECT * FROM company_types"
	var res []entities.CompanyTypes

	if err := comp.DB.Raw(query).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
