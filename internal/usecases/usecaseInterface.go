package usecases

import (
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
)

type CompanyUsecaseInterfaces interface {
	RegisterCompany(entities.CompanyResUsecase, string) (entities.CompanyResUsecase, error)
	AttachRolewithPremission(entities.CompanyRoles) error
	AddMember(entities.CompanyMembers) error
	GetRolesWithPermissions(compID string) ([]entities.CompanyRoles, error)
	GetPermissions() ([]entities.Permissions, error)
	GetCompanyTypes() ([]entities.CompanyTypes, error)
	AddCompanyType(entities.CompanyTypes) error
	AddPermissions(entities.Permissions) error
	GetCompanyDetails(string) (entities.ComapnyDetailsUsecase, error)
	GetCompanyMembers(string) ([]entities.GetCompanyEmployeesUsecase, error)
	AddMemberStatueses(string) error
	GetAverageSalaryperRole(string) ([]entities.AverageSalaryperRoleUsecase, error)
	RaiseProblem(entities.Problems) error
	GetProblems(string) ([]entities.Problems, error)
	InsertVisitors(string, string) error
	GetVisitorsWithinTimeframe(string, time.Time, time.Time) ([]entities.Visitors, error)
	GetVisitors(string) ([]entities.Visitors, error)
	GetProfileViews(string, time.Time, time.Time) (int, error)
	SalaryIncrementofEmployee(string, string, int) error
	SalaryIncrementofRole(string, uint, int) error
	LogintoCompany(string, string) (entities.LogintoCompanyUsecase, error)
	GetEmployeeLeaderBoard(string) ([]entities.GetEmployeeLeaderBoardUsecase, error)
	IsOwner(string, string) (bool, error)
	GetPermission(uint) (string, error)
	IsEmployeeExists(string, string) (bool, error)
	InsertIntoClients(entities.Clients) error
	AttachClientwithProject(entities.Clients, string, uint) error
	GetPastProjects(string) ([]entities.GetPastProjectsUsecase, error)
	GetClients(string) ([]entities.GetClientsUsecase, error)
	GetRevenuesGenerated(string) ([]entities.GetRevenueGeneratedUsecase, error)
	UpdateRevenueStatus(entities.UpdateRevenueStatusUsecase) error
	UpdateCompanyPolicies(entities.CompanyPolicies) error
	UpdatePayRollofEmployee(entities.PayRoll) error
	AssignProblemToEmployee(string, uint) error
	ResolveProblem(uint, string) error
	ApplyforLeave(entities.Leaves) error
	GetAppliedLeaves(string) ([]entities.Leaves, error)
	GrantLeave(uint, bool) error
	GetLeaves(string) ([]entities.Leaves, error)
	GetClientID(string) (string, error)
	PostJob(entities.Address, entities.Jobs) error
	ApplyJob(entities.JobApplications)(error)
	GetJobsofCompany(string)([]entities.Jobs,error)
	GetApplicationsforJob(string)([]entities.JobApplications,error)
	ShortlistApplications(string)(error)
	ScheduleInterviews(entities.ScheduledInterviews)(error)
	GetScheduledInterviews(string)([]entities.ScheduledInterviews,error)
	GetDetialsodApplicationbyID(string) (entities.JobApplications, error) 
	GetScheduledInterviewsofUser(string)([]entities.ScheduledInterviews,error)
	RescheduleInterview(entities.ScheduledInterviews)(error)
	GetShortlistedApplications(string)([]entities.JobApplications,error)
	GetJobs(string,string) ([]entities.Jobs, error)
	GetJobApplicationsofUser(string)([]entities.GetJobApplicationsofUserUsecase,error)
	GetAssignedProblems(string,string) ([]entities.Problems, error)
	DropCompany(string)(error)
	EditCompanyDetails(entities.Credentials)(error)
	EditCompanyEmployees(entities.CompanyMembers)(error)
	DeleteProblem(uint)(error)
	EditProblem(entities.Problems)(error)
	DropClient(entities.Clients)(error)
	UpdateCompanypolicies(entities.CompanyPolicies)(error)
	DeleteJob(string)(error)
	UpdateJob(entities.Jobs)(error)
	TerminateEmployee(string,string)(error)
}
