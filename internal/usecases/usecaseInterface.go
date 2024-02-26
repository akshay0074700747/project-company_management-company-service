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
}
