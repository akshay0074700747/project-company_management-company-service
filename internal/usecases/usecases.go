package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/akshay0074700747/project-company_management-company-service/helpers"
	"github.com/akshay0074700747/project-company_management-company-service/internal/adapters"
	"github.com/minio/minio-go/v7"
)

type CompanyUseCases struct {
	Adapter adapters.CompanyAdapterInterfaces
}

func NewCompanyUseCases(adapter adapters.CompanyAdapterInterfaces) *CompanyUseCases {
	return &CompanyUseCases{
		Adapter: adapter,
	}
}

func (comp *CompanyUseCases) RegisterCompany(req entities.CompanyResUsecase, ownerId string) (entities.CompanyResUsecase, error) {

	if req.CompCred.Name == "" {
		return entities.CompanyResUsecase{}, errors.New("the name cannot be empty")
	}

	if req.CompCred.CompanyUsername == "" {
		return entities.CompanyResUsecase{}, errors.New("the userNmae cannot be empty")
	}

	exists, err := comp.Adapter.IsCompanyUsernameExists(req.CompCred.CompanyUsername)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsCompanyUsernameExists adapter")
		return entities.CompanyResUsecase{}, err
	}

	if exists {
		return entities.CompanyResUsecase{}, errors.New("the username already exists")
	}

	req.CompCred.CompanyID = helpers.GenUuid()

	var resCreds entities.CompanyResUsecase

	resCreds.CompCred, resCreds.Email, resCreds.Phones, resCreds.Address, err = comp.Adapter.InsertCompanyCredentials(req.CompCred, req.Email, req.Phones, req.Address, ownerId)
	if err != nil {
		helpers.PrintErr(err, "error occured at InsertCompanyCredentials adapter")
		return resCreds, err
	}

	return resCreds, nil
}

func (comp *CompanyUseCases) AttachRolewithPremission(req entities.CompanyRoles) error {

	if err := comp.Adapter.AttachCompanyRoleAndPermissions(req); err != nil {
		helpers.PrintErr(err, "error occured at AttachCompanyRoleAndPermissions adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AddMember(req entities.CompanyMembers) error {

	exists, err := comp.Adapter.IsMemberExists(req.MemberID, req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsMemberExists adapter")
		return err
	}

	if exists {
		return errors.New("the user already exists")
	}

	exists, err = comp.Adapter.IsRoleIDExists(req.RoleID)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsRoleIDExists adapter")
		return err
	}

	if !exists {
		return errors.New("the roleid doesnt exist")
	}

	if err = comp.Adapter.AddMember(req); err != nil {
		helpers.PrintErr(err, "error occured at AddMember adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetRolesWithPermissions(compID string) ([]entities.CompanyRoles, error) {

	res, err := comp.Adapter.GetRoleWithPermissionIDs(compID)
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetRoleWithPermissionIDs adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetPermissions() ([]entities.Permissions, error) {

	res, err := comp.Adapter.GetPermissions()
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetPermissions adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetCompanyTypes() ([]entities.CompanyTypes, error) {

	res, err := comp.Adapter.GetCompanyTypes()
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetCompanyTypes adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) AddCompanyType(req entities.CompanyTypes) error {

	if err := comp.Adapter.AddCompanyType(req); err != nil {
		helpers.PrintErr(err, "eroor occured at AddCompanyType adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AddPermissions(req entities.Permissions) error {

	if err := comp.Adapter.AddPermissions(req); err != nil {
		helpers.PrintErr(err, "eroor occured at AddPermissions adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetCompanyDetails(companyID string) (entities.ComapnyDetailsUsecase, error) {

	cred, err := comp.Adapter.GetCompanyDetails(companyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetProjectDetails adapter")
		return entities.ComapnyDetailsUsecase{}, err
	}
	mem, err := comp.Adapter.GetNoofMembers(companyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetNoofMembers adapter")
		return entities.ComapnyDetailsUsecase{}, err
	}

	return entities.ComapnyDetailsUsecase{
		ComapanyID:      cred.CompanyID,
		CompanyUsername: cred.CompanyUsername,
		Name:            cred.Name,
		Members:         mem,
	}, nil
}

func (comp *CompanyUseCases) GetCompanyMembers(companyID string) ([]entities.GetCompanyEmployeesUsecase, error) {

	res, err := comp.Adapter.GetCompanyMembers(companyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetCompanyMembers usecase")
		return nil, err
	}

	return res, nil
}

func (company *CompanyUseCases) AddMemberStatueses(status string) error {

	if err := company.Adapter.AddMemberStatueses(status); err != nil {
		helpers.PrintErr(err, "error occured at AddMemberStatueses")
		return err
	}

	return nil
}

func (company *CompanyUseCases) GetAverageSalaryperRole(compID string) ([]entities.AverageSalaryperRoleUsecase, error) {

	res, err := company.Adapter.GetAverageSalaryperRole(compID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetAverageSalaryperRole adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) RaiseProblem(req entities.Problems) error {

	if err := comp.Adapter.RaiseProblem(req); err != nil {
		helpers.PrintErr(err, "error occured at RaiseProblem adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetProblems(companyID string) ([]entities.Problems, error) {

	problems, err := comp.Adapter.GetProblems(companyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetProblems adapter")
		return nil, err
	}

	return problems, nil
}

func (comp *CompanyUseCases) InsertVisitors(companyUsername, visitorID string) error {

	companyID, err := comp.Adapter.GetCompanyIDFromName(companyUsername)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetCompanyIDFromName adapter")
		return err
	}

	timeVisited := time.Now()

	if err = comp.Adapter.InsertVisitors(entities.Visitors{
		CompanyID:   companyID,
		VisitorID:   visitorID,
		VisitedTime: timeVisited,
	}); err != nil {
		helpers.PrintErr(err, "error happened at InsertVisitors adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetVisitorsWithinTimeframe(compID string, from time.Time, to time.Time) ([]entities.Visitors, error) {

	res, err := comp.Adapter.GetVisitorsWithinTimeframe(compID, from, to)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetVisitorsWithinTimeframe adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetVisitors(compID string) ([]entities.Visitors, error) {

	res, err := comp.Adapter.GetVisitors(compID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetVisitors adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetProfileViews(compID string, from time.Time, to time.Time) (int, error) {

	res, err := comp.Adapter.GetProfileViews(compID, from, to)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetProfileViews adapter")
		return 0, err
	}

	return res, nil
}

func (comp *CompanyUseCases) SalaryIncrementofEmployee(compID, userID string, increment int) error {

	if err := comp.Adapter.SalaryIncrementofEmployee(compID, userID, increment); err != nil {
		helpers.PrintErr(err, "error happened at SalaryIncrementofEmployee adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) SalaryIncrementofRole(compID string, roleId uint, increment int) error {

	if err := comp.Adapter.SalaryIncrementofRole(compID, roleId, increment); err != nil {
		helpers.PrintErr(err, "error happened at SalaryIncrementofRole adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) LogintoCompany(compUsername, userID string) (entities.LogintoCompanyUsecase, error) {

	fmt.Println(compUsername, userID, " heeereeeeeeeeee")
	compID, err := comp.Adapter.GetCompanyIDFromName(compUsername)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetCompanyIDFromName adapter")
		return entities.LogintoCompanyUsecase{}, err
	}
	fmt.Println(compID, "heeereeeeeeeeusernameeeeeeee")

	res, err := comp.Adapter.LogintoCompany(compID, userID)
	if err != nil {
		helpers.PrintErr(err, "error happened at LogintoCompany adapter")
		return entities.LogintoCompanyUsecase{}, err
	}

	res.CompanyID = compID

	return res, nil
}

func (comp *CompanyUseCases) GetEmployeeLeaderBoard(companyID string) ([]entities.GetEmployeeLeaderBoardUsecase, error) {

	res, err := comp.Adapter.GetEmployeeLeaderBoard(companyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetEmployeeLeaderBoard adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) IsOwner(user_id, company_id string) (bool, error) {

	res, err := comp.Adapter.IsOwner(company_id, user_id)

	fmt.Println(user_id, company_id, "from is owner")
	if err != nil {
		helpers.PrintErr(err, "error happened at IsOwner adapter")
		return false, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetPermission(id uint) (string, error) {

	permission, err := comp.Adapter.GetPermission(id)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetPermission")
		return "", err
	}

	return permission, nil
}

func (comp *CompanyUseCases) IsEmployeeExists(userID, compID string) (bool, error) {

	res, err := comp.Adapter.IsMemberExists(userID, compID)
	if err != nil {
		helpers.PrintErr(err, "errror happened at IsMemberExists adapter")
		return false, err
	}

	if !res {
		return false, nil
	}

	return true, nil
}

func (comp *CompanyUseCases) InsertIntoClients(req entities.Clients) error {

	if err := comp.Adapter.InsertIntoClients(req); err != nil {
		helpers.PrintErr(err, "error happened at InsertIntoClients")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AttachClientwithProject(req entities.Clients, projectid string, contract uint) error {

	if err := comp.Adapter.AttachClientwithProject(req, projectid, contract); err != nil {
		helpers.PrintErr(err, "error happened at AttachClientwithProject usecase")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetPastProjects(companyId string) ([]entities.GetPastProjectsUsecase, error) {

	res, err := comp.Adapter.GetPastProjects(companyId)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetPastProjects usecase")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetClients(compID string) ([]entities.GetClientsUsecase, error) {

	res, err := comp.Adapter.GetClients(compID)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetClients")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetRevenuesGenerated(compID string) ([]entities.GetRevenueGeneratedUsecase, error) {

	res, err := comp.Adapter.GetRevenuesGenerated(compID)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetRevenuesGenerated usecse")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) UpdateRevenueStatus(req entities.UpdateRevenueStatusUsecase) error {

	if err := comp.Adapter.UpdateRevenueStatus(req); err != nil {
		helpers.PrintErr(err, "error happened at UpdateRevenueStatus usecase")
		return err
	}

	return nil
}

func (company *CompanyUseCases) UpdateCompanyPolicies(req entities.CompanyPolicies) error {

	if err := company.Adapter.UpdateCompanyPolicies(req); err != nil {
		helpers.PrintErr(err, "eroro happened at UpdateCompanyPolicies usecase")
		return err
	}

	return nil
}

func (company *CompanyUseCases) UpdatePayRollofEmployee(req entities.PayRoll) error {

	if err := company.Adapter.UpdatePayRollofEmployee(req); err != nil {
		helpers.PrintErr(err, "eroror happeneded at UpdatePayRollofEmployee adapter")
		return err
	}

	return nil
}

func (company *CompanyUseCases) AssignProblemToEmployee(empID string, probID uint) error {

	if err := company.Adapter.AssignProblemToEmployee(empID, probID); err != nil {
		helpers.PrintErr(err, "eroro happened at AssignProblemToEmployee adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) ResolveProblem(probID uint, resolverID string) error {

	if err := comp.Adapter.ResolveProblem(probID, resolverID); err != nil {
		helpers.PrintErr(err, "error happended at ResolveProblem usecsae")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) ApplyforLeave(req entities.Leaves) error {

	if err := comp.Adapter.ApplyforLeave(req); err != nil {
		helpers.PrintErr(err, "error happened at ApplyforLeave adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetAppliedLeaves(compID string) ([]entities.Leaves, error) {

	res, err := comp.Adapter.GetAppliedLeaves(compID)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetAppliedLeaves adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GrantLeave(id uint, isAllowed bool) error {

	if err := comp.Adapter.GrantLeave(id, isAllowed); err != nil {
		helpers.PrintErr(err, "erorr happened at GrantLeave adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetLeaves(id string) ([]entities.Leaves, error) {

	res, err := comp.Adapter.GetLeaves(id)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetLeaves adpter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetClientID(projectID string) (string, error) {

	id, err := comp.Adapter.GetClientID(projectID)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetClientID adapter")
		return "", err
	}

	return id, nil
}

func (comp *CompanyUseCases) PostJob(address entities.Address, jobs entities.Jobs) error {

	jobs.JobID = helpers.GenUuid()

	if err := comp.Adapter.PostJob(address, jobs); err != nil {
		helpers.PrintErr(err, "error hapened at postjob adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) ApplyJob(req entities.JobApplications) error {

	req.ApplicationID = helpers.GenUuid()
	req.ResumeID = (helpers.GenUuid() + req.FileName)
	newReader := bytes.NewReader(req.Resume)

	if err := comp.Adapter.InsertResumetoMinio(context.TODO(), req.ResumeID, newReader, newReader.Size(), minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"applicationID": req.ApplicationID,
		},
	}); err != nil {
		helpers.PrintErr(err, "error happened at InsertResumetoMinio adapter")
		return err
	}

	if err := comp.Adapter.ApplyforJob(req); err != nil {
		helpers.PrintErr(err, "error happened at ApplyforJob usecase")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetJobsofCompany(compID string) ([]entities.Jobs, error) {

	res, err := comp.Adapter.GetJobsofCompany(compID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetJobsofCompany adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetApplicationsforJob(jobID string) ([]entities.JobApplications, error) {

	res, err := comp.Adapter.GetApplicationsforJob(jobID)
	if err != nil {
		helpers.PrintErr(err, "errror happened at GetApplicationsforJob adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) ShortlistApplications(applicationID string) error {

	if err := comp.Adapter.ShortlistApplications(applicationID); err != nil {
		helpers.PrintErr(err, "error happened at ShortlistApplications adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) ScheduleInterviews(req entities.ScheduledInterviews) error {

	if err := comp.Adapter.ScheduleInterviews(req); err != nil {
		helpers.PrintErr(err, "eror happened at ScheduleInterviews adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetScheduledInterviews(compID string) ([]entities.ScheduledInterviews, error) {

	res, err := comp.Adapter.GetScheduledInterviews(compID)
	if err != nil {
		helpers.PrintErr(err, "eroro happened at GetScheduledInterviews usecase")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetDetialsodApplicationbyID(appID string) (entities.JobApplications, error) {

	res, err := comp.Adapter.GetDetialsodApplicationbyID(appID)
	if err != nil {
		helpers.PrintErr(err, "erorr happened at GetDetialsodApplicationbyID adapter")
		return entities.JobApplications{}, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetScheduledInterviewsofUser(userID string) ([]entities.ScheduledInterviews, error) {

	res, err := comp.Adapter.GetScheduledInterviews(userID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetScheduledInterviews adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) RescheduleInterview(req entities.ScheduledInterviews) error {

	if err := comp.Adapter.RescheduleInterview(req); err != nil {
		helpers.PrintErr(err, "error happeedne at RescheduleInterview adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetShortlistedApplications(jobID string) ([]entities.JobApplications, error) {

	res, err := comp.Adapter.GetShortlistedApplications(jobID)
	if err != nil {
		helpers.PrintErr(err, "error happeed at GetShortlistedApplications adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetJobs(companyID, role string) ([]entities.Jobs, error) {

	var cond = make(map[string]interface{})

	if companyID != "" {
		cond["company_id"] = companyID
	}
	if role != "" {
		cond["role"] = role
	}

	res, err := comp.Adapter.GetJobs(cond)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetJobs adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetJobApplicationsofUser(userID string) ([]entities.GetJobApplicationsofUserUsecase, error) {

	res, err := comp.Adapter.GetJobApplicationsofUser(userID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetJobApplicationsofUser adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetAssignedProblems(companyID, userID string) ([]entities.Problems, error) {

	res, err := comp.Adapter.GetAssignedProblems(companyID, userID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetAssignedProblems")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) DropCompany(compID string) error {

	if err := comp.Adapter.DropCompany(compID); err != nil {
		helpers.PrintErr(err, "erorr happened at DropCompany adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) EditCompanyDetails(req entities.Credentials) error {

	if err := comp.Adapter.EditCompanyDetails(req); err != nil {
		helpers.PrintErr(err, "error happenendd at EditCompanyDetails adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) EditCompanyEmployees(req entities.CompanyMembers) error {

	if err := comp.Adapter.EditCompanyEmployees(req); err != nil {
		helpers.PrintErr(err, "error happenendd at EditCompanyEmployees adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) DeleteProblem(id uint) error {

	if err := comp.Adapter.DeleteProblem(id); err != nil {
		helpers.PrintErr(err, "error happneed at DeleteProblem adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) EditProblem(req entities.Problems) error {

	if err := comp.Adapter.EditProblem(req); err != nil {
		helpers.PrintErr(err, "error happenend  at EditProblem adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) DropClient(req entities.Clients) error {

	if err := comp.Adapter.DropClient(req); err != nil {
		helpers.PrintErr(err, "error happened at DropClient adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) UpdateCompanypolicies(req entities.CompanyPolicies) error {

	if err := comp.Adapter.UpdateCompanypolicies(req); err != nil {
		helpers.PrintErr(err, "erorro happened at UpdateCompanypolicies adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) DeleteJob(jobID string) error {

	if err := comp.Adapter.DeleteJob(jobID); err != nil {
		helpers.PrintErr(err, "error happened at DeleteJob adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) UpdateJob(req entities.Jobs) error {

	if err := comp.Adapter.UpdateJob(req); err != nil {
		helpers.PrintErr(err, "error happened at UpdateJob adapter")
		return err
	}

	return nil
}
