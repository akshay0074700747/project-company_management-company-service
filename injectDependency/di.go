package injectdependency

import (
	"github.com/akshay0074700747/project-company_management-company-service/config"
	"github.com/akshay0074700747/project-company_management-company-service/db"
	"github.com/akshay0074700747/project-company_management-company-service/internal/adapters"
	"github.com/akshay0074700747/project-company_management-company-service/internal/services"
	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
)

func Initialize(cfg config.Config) *services.CompanyEngine {

	dbb := db.ConnectDB(cfg)
	minio := db.ConnectMinio(cfg)
	adapter := adapters.NewCompanyAdapter(dbb,minio)
	usecase := usecases.NewCompanyUseCases(adapter)
	server := services.NewProjectServiceServer(usecase, ":50001",":50002")
	go server.StartConsuming()
	
	return services.NewCompanyEngine(server)
}
