package adapters

import "github.com/akshay0074700747/project-company_management-company-service/entities"

type CompanyAdapterInterfaces interface {
	InsertCompanyCredentials(entities.Credentials) (entities.Credentials, error)
	InsertEmail([]entities.CompanyEmail) ([]entities.CompanyEmail, error)
	InsertPhone([]entities.CompanyPhone) ([]entities.CompanyPhone, error)
	InsertAddress(entities.CompanyAddress) (entities.CompanyAddress, error)
	IsCompanyUsernameExists(string) (bool, error)
	AttachCompanyRoleAndPermissions(entities.CompanyRoles) error
	AddMember(entities.CompanyMembers) error
	IsMemberExists(string) (bool, error)
	IsRoleIDExists(uint) (bool, error)
	GetRoleWithPermissionIDs(string) ([]entities.CompanyRoles,error)
	GetPermissions()([]entities.Permissions,error)
	GetCompanyTypes()([]entities.CompanyTypes,error)
}
