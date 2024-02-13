package usecases

import "github.com/akshay0074700747/project-company_management-company-service/internal/adapters"



type CompanyUseCases struct {
	Adapter adapters.CompanyAdapterInterfaces
}

func NewCompanyUseCases(adapter adapters.CompanyAdapterInterfaces) *CompanyUseCases {
	return &CompanyUseCases{
		Adapter: adapter,
	}
}
