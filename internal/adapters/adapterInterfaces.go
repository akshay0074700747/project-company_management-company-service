package adapters

import (
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
)

type CompanyAdapterInterfaces interface {
	InsertCompanyCredentials(entities.Credentials, []string, []string, entities.CompanyAddress) (entities.Credentials, []string, []string, entities.CompanyAddress, error)
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
}
