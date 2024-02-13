package services

import (
	"context"

	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/companypb"
)

type CompanyServiceServer struct {
	Usecase usecases.CompanyUsecaseInterfaces
	companypb.UnimplementedCompanyServiceServer
}

func NewProjectServiceServer(usecase usecases.CompanyUsecaseInterfaces) *CompanyServiceServer {
	return &CompanyServiceServer{
		Usecase: usecase,
	}
}

func (auth *CompanyServiceServer) RegisterCompany(ctx context.Context, req *companypb.RegisterCompanyRequest) (*companypb.CompanyResponce, error) {
	
}
