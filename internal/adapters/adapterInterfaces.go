package adapters

import (
	"context"
	"io"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/minio/minio-go/v7"
)

type CompanyAdapterInterfaces interface {
	InsertCompanyCredentials(entities.Credentials, []string, []string, entities.CompanyAddress, string) (entities.Credentials, []string, []string, entities.CompanyAddress, error)
	InsertEmail([]entities.CompanyEmail) ([]entities.CompanyEmail, error)
	InsertPhone([]entities.CompanyPhone) ([]entities.CompanyPhone, error)
	InsertAddress(entities.CompanyAddress) (entities.CompanyAddress, error)
	IsCompanyUsernameExists(string) (bool, error)
	AttachCompanyRoleAndPermissions(entities.CompanyRoles) error
	AddMember(entities.CompanyMembers) error
	IsMemberExists(string, string) (bool, error)
	IsRoleIDExists(uint) (bool, error)
	GetRoleWithPermissionIDs(string) ([]entities.CompanyRoles, error)
	GetPermissions() ([]entities.Permissions, error)
	GetCompanyTypes() ([]entities.CompanyTypes, error)
	AddOwner(string, string) error
	AddCompanyType(entities.CompanyTypes) error
	AddPermissions(entities.Permissions) error
	GetNoofMembers(string) (uint, error)
	GetCompanyDetails(string) (entities.Credentials, error)
	AddMemberStatueses(string) error
	GetAverageSalaryperRole(string) ([]entities.AverageSalaryperRoleUsecase, error)
	RaiseProblem(entities.Problems) error
	GetProblems(string) ([]entities.Problems, error)
	GetCompanyIDFromName(string) (string, error)
	InsertVisitors(entities.Visitors) error
	GetVisitorsWithinTimeframe(string, time.Time, time.Time) ([]entities.Visitors, error)
	GetVisitors(string) ([]entities.Visitors, error)
	GetProfileViews(string, time.Time, time.Time) (int, error)
	SalaryIncrementofEmployee(string, string, int) error
	SalaryIncrementofRole(string, uint, int) error
	LogintoCompany(string, string) (entities.LogintoCompanyUsecase, error)
	GetEmployeeLeaderBoard(string) ([]entities.GetEmployeeLeaderBoardUsecase, error)
	GetCompanyMembers(string) ([]entities.GetCompanyEmployeesUsecase, error)
	IsOwner(string, string) (bool, error)
	GetPermission(uint) (string, error)
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
	InsertResumetoMinio(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error
	ApplyforJob(entities.JobApplications) error
	GetJobsofCompany(string)([]entities.Jobs,error)
	GetApplicationsforJob(string)([]entities.JobApplications,error)
	ShortlistApplications(string)(error)
	ScheduleInterviews(entities.ScheduledInterviews)(error)
	GetScheduledInterviews(string)([]entities.ScheduledInterviews,error)
	GetDetialsodApplicationbyID(string)(entities.JobApplications,error)
	GetScheduledInterviewsofUser(string)([]entities.ScheduledInterviews,error)
	RescheduleInterview(entities.ScheduledInterviews)(error)
	GetShortlistedApplications(string)([]entities.JobApplications,error)
	GetJobs(map[string]interface{}) ([]entities.Jobs, error)
	GetJobApplicationsofUser(string)([]entities.GetJobApplicationsofUserUsecase,error)
	GetAssignedProblems(string,string)([]entities.Problems,error)
	DropCompany(string)(error)
	EditCompanyDetails(entities.Credentials)(error)
	EditCompanyEmployees(entities.CompanyMembers)(error)
	DeleteProblem(uint)(error)
	EditProblem(entities.Problems)(error)
	DropClient(entities.Clients)(error)
	UpdateCompanypolicies(entities.CompanyPolicies)(error)
	DeleteJob(string)(error)
	UpdateJob(entities.Jobs)(error)
}