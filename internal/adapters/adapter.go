package adapters

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type CompanyAdapter struct {
	DB      *gorm.DB
	MinioDB *minio.Client
}

func NewCompanyAdapter(db *gorm.DB, minioDB *minio.Client) *CompanyAdapter {
	return &CompanyAdapter{
		DB:      db,
		MinioDB: minioDB,
	}
}

func (company *CompanyAdapter) InsertCompanyCredentials(req entities.Credentials, emailss, phoness []string, addr entities.CompanyAddress, ownerID string) (entities.Credentials, []string, []string, entities.CompanyAddress, error) {

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
	fmt.Println(req.CompanyID, ownerID)
	if err := tx.Exec("INSERT INTO owners (company_id,owner_id) VALUES($1,$2)", req.CompanyID, ownerID).Error; err != nil {
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

	query := "SELECT role_id,AVG(salary) AS salary FROM company_members WHERE company_id = $1 GROUP BY role_id"
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

	fmt.Println(companyUsername)
	query := "SELECT company_id FROM credentials WHERE company_username = $1"
	var res string

	tx := comp.DB.Raw(query, companyUsername).Scan(&res)
	if tx.Error != nil {
		return "", tx.Error
	}

	fmt.Println(res)

	if tx.RowsAffected == 0 {
		return "", errors.New("the companyUsername is not valid")
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

	var result entities.LogintoCompanyUsecase
	var res entities.LogintoCompanyUsecaseTmp
	query := "SELECT m.company_id,r.role_id FROM company_members m INNER JOIN company_roles r ON r.company_id = $1 AND r.id = m.role_id WHERE m.company_id = $1 AND m.member_id = $2"
	tx := comp.DB.Raw(query, compID, userID).Scan(&res)
	if tx.Error != nil {
		return entities.LogintoCompanyUsecase{}, tx.Error
	}

	query = "SELECT permission_id FROM company_roles WHERE company_id = $1 AND role_id = $2"
	var permID uint
	if err := comp.DB.Raw(query, res.CompanyID, res.RoleID).Scan(&permID).Error; err != nil {
		return entities.LogintoCompanyUsecase{}, err
	}

	var perm string
	query = "SELECT permission FROM permissions WHERE id = $1"
	if err := comp.DB.Raw(query, permID).Scan(&perm).Error; err != nil {
		return entities.LogintoCompanyUsecase{}, err
	}

	result.CompanyID = res.CompanyID
	result.RoleID = res.RoleID
	result.Permisssion = perm

	fmt.Println(result, "----------------")

	return result, nil
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

	fmt.Println(companyID, "---company id")
	query := "SELECT m.member_id AS user_id,r.role_id,p.permission FROM company_members m INNER JOIN company_roles r ON r.id = m.role_id AND r.company_id = $1 INNER JOIN permissions p ON p.id = r.permission_id WHERE m.company_id = $1"

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	fmt.Println(res, "---result")

	return res, nil
}

func (comp *CompanyAdapter) IsOwner(user_id, company_id string) (bool, error) {

	query := "SELECT * FROM owners WHERE company_id = $1 AND owner_id = $2"
	tx := comp.DB.Exec(query, company_id, user_id)
	fmt.Println(company_id, user_id)
	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (comp *CompanyAdapter) GetPermission(id uint) (string, error) {

	query := "SELECT permission FROM permissions WHERE id = $1"
	var res string

	if err := comp.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return "", err
	}

	return res, nil
}

func (comp *CompanyAdapter) InsertIntoClients(req entities.Clients) error {

	query := "INSERT INTO clients (client_id,company_id) VALUES($1,$2)"
	if err := comp.DB.Exec(query, req.ClientID, req.CompanyID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) AttachClientwithProject(req entities.Clients, projectID string, contract uint) error {

	query := "INSERT INTO clients_with_projects (client_id,project_id,contract) VALUES((SELECT id FROM clients c WHERE c.client_id = $1 AND c.company_id = $2),$3,$4)"
	tx := comp.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Exec(query, req.ClientID, req.CompanyID, projectID, contract).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func (comp *CompanyAdapter) GetPastProjects(companyID string) ([]entities.GetPastProjectsUsecase, error) {

	var res []entities.GetPastProjectsUsecase
	query := "SELECT c.client_id,p.project_id FROM clients c INNER JOIN clients_with_projects p ON c.id = p.client_id AND p.is_recieved = true WHERE c.company_id = $1"

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetClients(compId string) ([]entities.GetClientsUsecase, error) {

	var res []entities.GetClientsUsecase
	query := "SELECT c.client_id,p.project_id FROM clients c LEFT JOIN clients_with_projects p ON c.id = p.client_id WHERE c.company_id = $1"

	if err := comp.DB.Raw(query, compId).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetRevenuesGenerated(compID string) ([]entities.GetRevenueGeneratedUsecase, error) {

	query := "SELECT p.project_id,p.contract AS revenue,c.client_id FROM clients_with_projects p INNER JOIN clients c ON c.id = p.client_id AND c.company_id = $1 WHERE p.is_recieved = true"
	var res []entities.GetRevenueGeneratedUsecase

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) UpdateRevenueStatus(req entities.UpdateRevenueStatusUsecase) error {

	query := "UPDATE clients_with_projects SET is_recieved = $1 WHERE project_id = $2 AND client_id = (SELECT id FROM clients WHERE client_id = $3 AND company_id = $4)"
	if err := comp.DB.Exec(query, req.IsRecieved, req.ProjectID, req.ClientID, req.CompanyID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) UpdateCompanyPolicies(req entities.CompanyPolicies) error {

	selectQuery := "SELECT * FROM company_policies WHERE company_id = $1"
	tx := comp.DB.Raw(selectQuery, req.CompanyID)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {

		query := "INSERT INTO company_policies (company_id,maxleaves_per_month,pay_day,working_hours_perday) VALUES($1,$2,$3,$4)"
		if err := comp.DB.Exec(query, req.CompanyID, req.MaxleavesPerMonth, req.PayDay, req.WorkingHoursPerday).Error; err != nil {
			return err
		}

	} else {

		query := "UPDATE company_policies SET maxleaves_per_month = $1,pay_day = $2,working_hours_perday = $3 WHERE company_id = $4"
		if err := comp.DB.Exec(query, req.MaxleavesPerMonth, req.PayDay, req.WorkingHoursPerday, req.CompanyID).Error; err != nil {
			return err
		}
	}

	return nil
}

func (comp *CompanyAdapter) UpdatePayRollofEmployee(req entities.PayRoll) error {

	query := "INSERT INTO pay_rolls (company_id,employee_id,is_payed,transaction_id) VALUES($1,$2,$3,$4)"
	if err := comp.DB.Exec(query, req.CompanyID, req.EmployeeID, req.IsPayed, req.TransactionID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) AssignProblemToEmployee(empID string, probID uint) error {

	query := "UPDATE problems SET assigned_employee_id = $1 WHERE id = $2"
	if err := comp.DB.Exec(query, empID, probID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) ResolveProblem(probID uint, resolverID string) error {

	query := "UPDATE problems SET is_resolved = true AND assigned_employee_id = $1 WHERE id = $2"
	if err := comp.DB.Exec(query, resolverID, probID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) ApplyforLeave(req entities.Leaves) error {

	query := "INSERT INTO leaves (employee_id,company_id,description,date) VALUES($1,$2,$3,$4)"
	if err := comp.DB.Exec(query, req.EmployeeID, req.CompanyID, req.Description, req.Date).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetAppliedLeaves(compID string) ([]entities.Leaves, error) {

	query := "SELECT * FROM leaves WHERE company_id = $1"
	var res []entities.Leaves

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GrantLeave(id uint, is_allowed bool) error {

	query := "UPDATE leaves SET is_allowed = $1 WHERE id = $2"
	if err := comp.DB.Exec(query, is_allowed, id).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetLeaves(compID string) ([]entities.Leaves, error) {

	query := "SELECT * FROM leaves WHERE company_id = $1 AND is_allowed = true"
	var res []entities.Leaves

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetClientID(projectID string) (string, error) {

	query := "SELECT client_id FROM clients WHERE id = (SELECT client_id FROM clients_with_projects WHERE project_id = $1)"
	var clientID string

	if err := comp.DB.Raw(query, projectID).Scan(&clientID).Error; err != nil {
		return "", err
	}

	return clientID, nil
}

func (comp *CompanyAdapter) PostJob(address entities.Address, jobs entities.Jobs) error {

	tx := comp.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := "INSERT INTO addresses (street_no,street_name,pin_no,district,state,nation) VALUES($1,$2,$3,$4,$5,$6) RETURNING id"
	var id uint
	if err := tx.Raw(query, address.StreetNo, address.StreetName, address.PinNo, address.District, address.State, address.Nation).Scan(&id).Error; err != nil {
		tx.Rollback()
		return err
	}

	query = "INSERT INTO jobs (job_id,company_id,role,vacancy,description,min_experiance,min_expected_ctc,max_expected_ctc,is_remote,address_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	if err := tx.Exec(query, jobs.JobID, jobs.CompanyID, jobs.Role, jobs.Vacancy, jobs.Description, jobs.MinExperiance, jobs.MinExpectedCTC, jobs.MaxExpectedCTC, jobs.IsRemote, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (company *CompanyAdapter) InsertResumetoMinio(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error {

	_, err := company.MinioDB.PutObject(ctx, "resume-storage-bucket", fileName, reader, objectSize, opts)
	if err != nil {
		return err
	}

	return nil
}

func (company *CompanyAdapter) ApplyforJob(req entities.JobApplications) error {

	tx := company.DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("panic occured",r)
	//		tx.Rollback()
	//	}
	// }()

	// query := "INSERT INTO addresses (street_name,street_no,district,state,nation,pin_no) VALUES($1,$2,$3,$4,$5,$6) RETURNING id"
	// var id uint

	if err := tx.Create(&req.AddressofApplicant).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	req.AddressID = req.AddressofApplicant.ID

	// query = "INSERT INTO job_applications (application_id,user_id,name,email,mobile,highest_education,nationality,experiance,current_ctc,resume_id,address_id,job_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)"

	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	tx.Commit()

	return nil
}

func (comp *CompanyAdapter) GetJobsofCompany(companyID string) ([]entities.Jobs, error) {

	query := "SELECT * FROM jobs WHERE company_id = $1"
	var res []entities.Jobs

	if err := comp.DB.Raw(query, companyID).Scan(&res).Error; err != nil {
		return nil, err
	}

	query = "SELECT COUNT(*) FROM job_applications WHERE job_id = $1 GROUP BY job_id"
	for i := range res {
		if err := comp.DB.Raw(query, res[i].JobID).Scan(&res[i].TotalPersonsApplied).Error; err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (comp *CompanyAdapter) GetApplicationsforJob(jobID string) ([]entities.JobApplications, error) {

	query := "SELECT * FROM job_applications WHERE job_id = $1"
	var res []entities.JobApplications

	if err := comp.DB.Raw(query, jobID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) ShortlistApplications(applicationID string) error {

	query := "UPDATE job_applications SET is_shortlisted = true WHERE application_id = $1"
	if err := comp.DB.Exec(query, applicationID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) ScheduleInterviews(req entities.ScheduledInterviews) error {

	query := "INSERT INTO scheduled_interviews (application_id,date,description,time) VALUES($1,$2,$3,$4)"

	if err := comp.DB.Exec(query, req.ApplicationID, req.Date, req.Description, req.Time).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetScheduledInterviews(compID string) ([]entities.ScheduledInterviews, error) {

	query := "SELECT * FROM scheduled_interviews WHERE application_id IN (SELECT application_id FROM job_applications WHERE job_id IN (SELECT job_id FROM jobs WHERE company_id = $1))"
	var res []entities.ScheduledInterviews

	if err := comp.DB.Raw(query, compID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetDetialsodApplicationbyID(appID string) (entities.JobApplications, error) {

	query := "SELECT * FROM job_applications WHERE application_id = $1"
	var res entities.JobApplications

	if err := comp.DB.Raw(query, appID).Scan(&res).Error; err != nil {
		return entities.JobApplications{}, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetScheduledInterviewsofUser(userID string) ([]entities.ScheduledInterviews, error) {

	query := "SELECT * FROM scheduled_interviews WHERE application_id IN (SELECT application_id FROM job_applications WHERE user_id = $1)"
	var res []entities.ScheduledInterviews

	if err := comp.DB.Raw(query, userID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) RescheduleInterview(req entities.ScheduledInterviews) error {

	query := "UPDATE scheduled_interviews SET date = $1,description = $2,time = $3 WHERE application_id = $4"

	if err := comp.DB.Raw(query, req.Date, req.Description, req.Time, req.ApplicationID).Error; err != nil {
		return err
	}

	return nil
}

func (comp *CompanyAdapter) GetShortlistedApplications(jobID string) ([]entities.JobApplications, error) {

	query := "SELECT * FROM job_applications WHERE job_id = $1 AND is_shortlisted = true"
	var res []entities.JobApplications

	if err := comp.DB.Raw(query, jobID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetJobs(conditions map[string]interface{}) ([]entities.Jobs, error) {

	var res []entities.Jobs
	if err := comp.DB.Model(&entities.Jobs{}).Where(conditions).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetJobApplicationsofUser(userID string) ([]entities.GetJobApplicationsofUserUsecase, error) {

	query := "SELECT j.company_id,j.role,ja.application_id,ja.is_verified FROM jobs j INNER JOIN job_applications ja ON j.job_id = ja.job_id AND ja.user_id = $1"
	var res []entities.GetJobApplicationsofUserUsecase

	if err := comp.DB.Raw(query, userID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (comp *CompanyAdapter) GetAssignedProblems(companyID, userID string) ([]entities.Problems, error) {

	query := "SELECT * FROM problems WHERE company_id = $1 AND assigned_employee_id = $2 AND is_resolved = false"
	var res []entities.Problems

	if err := comp.DB.Raw(query, companyID, userID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
