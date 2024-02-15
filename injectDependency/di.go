package injectdependency

import (
	"github.com/akshay0074700747/project-company_management-company-service/config"
	"github.com/akshay0074700747/project-company_management-company-service/db"
	"github.com/akshay0074700747/project-company_management-company-service/internal/adapters"
	"github.com/akshay0074700747/project-company_management-company-service/internal/services"
	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
)

func Initialize(cfg config.Config) *services.CompanyEngine {

	db := db.ConnectDB(cfg)
	adapter := adapters.NewCompanyAdapter(db)
	usecase := usecases.NewCompanyUseCases(adapter)
	server := services.NewProjectServiceServer(usecase, ":50001")

	return services.NewCompanyEngine(server)
}
