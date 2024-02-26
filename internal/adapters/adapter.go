package adapters

import (
	"fmt"
	"time"

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

func (company *CompanyAdapter) InsertCompanyCredentials(req entities.Credentials, emailss, phoness []string, addr entities.CompanyAddress) (entities.Credentials, []string, []string, entities.CompanyAddress, error) {

	tx := company.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := "INSERT INTO credentials (company_id,company_username,name,type_id) VALUES($1,$2,$3,$4) RETURNING company_id,company_username,name,type_id"
	var res entities.Credentials

	if err := tx.Raw(query, req.CompanyID, req.CompanyUsername, req.Name, req.TypeID).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, []string{}, []string{}, entities.CompanyAddress{}, err
	}

	var emails []entities.CompanyEmail
	var phones []entities.CompanyPhone

	for _, v := range emailss {
		emails = append(emails, entities.CompanyEmail{
			Email:     v,
			CompanyID: res.CompanyID,
		})
	}

	for _, v := range phoness {
		phones = append(phones, entities.CompanyPhone{
			Phone:     v,
			CompanyID: res.CompanyID,
		})
	}

	if err := tx.Create(&emails).Error; err != nil {
		tx.Rollback()
		fmt.Println("----here was the error")
		return res, []string{}, []string{}, entities.CompanyAddress{}, err
	}

	if err := tx.Create(&phones).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, emailss, []string{}, entities.CompanyAddress{}, err
	}

	if err := tx.Create(&addr).Error; err != nil {
		tx.Rollback()
		return res, emailss, phoness, entities.CompanyAddress{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, emailss, phoness, addr, err
	}
	return res, emailss, phoness, addr, nil
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

	query := "INSERT INTO company_members (company_id,role_id,member_id,status_id,salary) VALUES($1,$2,$3,(SELECT id FROM member_statuses WHERE status = 'LIVE'),$4)"

	if err := comp.DB.Exec(query, req.CompanyID, req.RoleID, req.MemberID, req.Salary).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) IsMemberExists(id string, compID string) (bool, error) {

	query := "SELECT * FROM company_members WHERE member_id = $1 AND company_id = $2"

	res := comp.DB.Exec(query, id, compID)
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

func (comp *CompanyAdapter) AddOwner(compID, ownerID string) error {

	query := "INSERT INTO owners (company_id,owner_id) VALUES($1,$2)"

	if err := comp.DB.Exec(query, compID, ownerID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) AddCompanyType(req entities.CompanyTypes) error {

	query := "INSERT INTO company_types (type) VALUES($1)"

	if err := comp.DB.Exec(query, req.Type).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) AddPermissions(req entities.Permissions) error {

	query := "INSERT INTO permissions (permission) VALUES($1)"

	if err := comp.DB.Exec(query, req.Permission).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetNoofMembers(compID string) (uint, error) {

	query := "SELECT COUNT(*) FROM company_members WHERE company_id = $1 GROUP BY company_id"
	var res uint

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetCompanyDetails(compID string) (entities.Credentials, error) {

	query := "SELECT * FROM credentials WHERE company_id = $1"
	var res entities.Credentials

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return entities.Credentials{}, err
	}

	return res, nil
}

func (project *CompanyAdapter) AddMemberStatueses(status string) error {

	query := "INSERT INTO member_statuses (status) VALUES($1)"
	if err := project.DB.Exec(query, status).Error; err != nil {
		return err
	}

	return nil
}

func (company *CompanyAdapter) GetAverageSalaryperRole(compID string) ([]entities.AverageSalaryperRoleUsecase, error) {

	query := "SELECT role_id,AVG(salary) AS salary FROM members WHERE company_id = $1 GROUP BY role_id"
	var res []entities.AverageSalaryperRoleUsecase

	if err := company.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) RaiseProblem(req entities.Problems) error {

	query := "INSERT INTO problems (problem,company_id,raised_by) VALUES($1,$2,$3)"
	if err := comp.DB.Exec(query, req.Problem, req.CompanyID, req.RaisedBy).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetProblems(comapnyID string) ([]entities.Problems, error) {

	query := "SELECT * FROM problems WHERE company_id = $1"
	var res []entities.Problems

	if err := comp.DB.Raw(query, comapnyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetCompanyIDFromName(companyUsername string) (string, error) {

	query := "SELECT company_id FROM credentials WHERE company_username = $1"
	var res string

	if err := comp.DB.Raw(query, companyUsername).Scan(&res).Error; err != nil {
		return "", err
	}

	return res, nil
}

func (comp *CompanyAdapter) InsertVisitors(req entities.Visitors) error {

	query := "INSERT INTO visitors (company_id,visitor_id,visited_time) VALUES($1,$2,$3)"
	if err := comp.DB.Exec(query, req.CompanyID, req.VisitorID, req.VisitedTime).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetVisitorsWithinTimeframe(companyID string, from time.Time, to time.Time) ([]entities.Visitors, error) {

	query := "SELECT visitor_id,visited_time FROM visitors WHERE company_id = $1 AND visited_time >= $2 AND visited_time <= $3 ORDER BY visited_time DESC"
	var res []entities.Visitors

	if err := comp.DB.Raw(query, companyID, from, to).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetVisitors(companyID string) ([]entities.Visitors, error) {

	quer := "SELECT visitor_id,visited_time FROM visitors WHERE company_id = $1 ORDER BY visited_time DESC"
	var res []entities.Visitors

	if err := comp.DB.Raw(quer, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetProfileViews(compID string, from time.Time, to time.Time) (int, error) {

	query := "SELECT COUNT(*) AS views FROM visitors WHERE company_id = $1 AND visited_time >= $2 AND visited_time <= $3 GROUP BY company_id"
	var views int

	if err := comp.DB.Raw(query, compID, from, to).Scan(&views).Error; err != nil {
		return 0, err
	}

	return views, nil
}

func (comp *CompanyAdapter) SalaryIncrementofEmployee(companyID, userID string, increment int) error {

	query := "UPDATE company_members SET salary = salary + $1 WHERE company_id = $2 AND member_id = $3"

	if err := comp.DB.Exec(query, increment, companyID, userID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) SalaryIncrementofRole(compID string, roleID uint, increment int) error {

	query := "UPDATE company_members SET salary = salary + $1 WHERE company_id  = $2 AND role_id = $3"

	if err := comp.DB.Exec(query, increment, compID, roleID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) LogintoCompany(compID, userID string) (entities.LogintoCompanyUsecase, error) {

	var res entities.LogintoCompanyUsecase
	query := "SELECT m.company_id,p.permission,r.role_id FROM company_members m INNER JOIN company_roles r ON r.company_id = $1 AND r.id = m.role_id INNER JOIN permissions p ON p.id = r.permission_id WHERE m.company_id = $1 AND m.member_id = $2"
	if err := comp.DB.Raw(query, compID, userID).Scan(&res).Error; err != nil {
		return entities.LogintoCompanyUsecase{}, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetEmployeeLeaderBoard(companyID string) ([]entities.GetEmployeeLeaderBoardUsecase, error) {

	query := "SELECT m.member_id AS employee_id,m.salary,r.role_id FROM company_members m INNER JOIN company_roles r ON r.id = m.role_id AND r.company_id = $1 WHERE m.company_id = $1 ORDER BY m.salary DESC"
	var res []entities.GetEmployeeLeaderBoardUsecase

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetCompanyMembers(companyID string) ([]entities.GetCompanyEmployeesUsecase, error) {

	var res []entities.GetCompanyEmployeesUsecase
	query := "SELECT m.member_id AS user_id,r.role_id,p.permission FROM company_members m INNER JOIN company_roles r ON r.id = m.role_id AND r.company_id = $1 INNER JOIN permissions p ON p.id = r.permission_id WHERE company_id = $1"

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
