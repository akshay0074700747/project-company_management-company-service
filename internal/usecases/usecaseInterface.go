package usecases

import "github.com/akshay0074700747/project-company_management-company-service/entities"

type CompanyUsecaseInterfaces interface {
	RegisterCompany(entities.CompanyResUsecase) (entities.CompanyResUsecase, error)
	AttachRolewithPremission(entities.CompanyRoles) error
	AddMember(entities.CompanyMembers) error
	GetRolesWithPermissions(compID string) ([]entities.CompanyRoles, error)
	GetPermissions() ([]entities.Permissions, error)
	GetCompanyTypes() ([]entities.CompanyTypes, error)
}
