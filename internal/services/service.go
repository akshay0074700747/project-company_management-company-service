package services

import (
	"context"
	"errors"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/akshay0074700747/project-company_management-company-service/helpers"
	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/companypb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CompanyServiceServer struct {
	UserConn userpb.UserServiceClient
	Usecase  usecases.CompanyUsecaseInterfaces
	companypb.UnimplementedCompanyServiceServer
}

func NewProjectServiceServer(usecase usecases.CompanyUsecaseInterfaces, addr string) *CompanyServiceServer {
	userRes, _ := helpers.DialGrpc(addr)
	return &CompanyServiceServer{
		Usecase:  usecase,
		UserConn: userpb.NewUserServiceClient(userRes),
	}
}

func (auth *CompanyServiceServer) RegisterCompany(ctx context.Context, req *companypb.RegisterCompanyRequest) (*companypb.CompanyResponce, error) {

	res, err := auth.Usecase.RegisterCompany(entities.CompanyResUsecase{
		CompCred: entities.Credentials{
			CompanyUsername: req.Companyusername,
			Name:            req.Name,
			TypeID:          uint(req.TypeID),
		},
		Email:  req.Emails,
		Phones: req.Mobiles,
		Address: entities.CompanyAddress{
			StreetNo:   uint(req.Address.StreetNo),
			StreetName: req.Address.Street,
			PinNo:      uint(req.Address.PinNo),
			District:   req.Address.District,
			State:      req.Address.State,
			Nation:     req.Address.Nation,
		},
	})
	if err != nil {
		helpers.PrintErr(err, "error happened at RegisterCompany usecase")
		return nil, err
	}

	return &companypb.CompanyResponce{
		Name:            res.CompCred.Name,
		Companyusername: res.CompCred.CompanyUsername,
		CompanyID:       res.Address.CompanyID,
		Emails:          res.Email,
		Mobiles:         res.Phones,
		Address: &companypb.Address{
			Street:   res.Address.StreetName,
			StreetNo: int32(res.Address.StreetNo),
			PinNo:    int32(res.Address.PinNo),
			District: res.Address.District,
			State:    res.Address.State,
			Nation:   res.Address.Nation,
		},
	}, nil
}

func (comp *CompanyServiceServer) AddEmployees(ctx context.Context, req *companypb.AddEmployeeReq) (*emptypb.Empty, error) {

	res, err := comp.UserConn.GetByEmail(ctx, &userpb.GetByEmailReq{
		Email: req.Email,
	})

	if err != nil {
		return nil, errors.New("the user service had come problems")
	}

	if err = comp.Usecase.AddMember(entities.CompanyMembers{
		CompanyID: req.CompanyID,
		RoleID:    uint(req.RoleID),
		MemberID:  res.UserID,
	}); err != nil {
		helpers.PrintErr(err, "error occured at AddMember usecase")
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetCompanyTypes(emp *emptypb.Empty, stream companypb.CompanyService_GetCompanyTypesServer) error {

	res, err := comp.Usecase.GetCompanyTypes()
	if err != nil {
		helpers.PrintErr(err, "eroro occured at GetCompanyTypes usecase")
		return err
	}
	for _, v := range res {
		if err = stream.Send(&companypb.GetCompanyTypesRes{
			ID:   uint32(v.ID),
			Name: v.Type,
		}); err != nil {
			helpers.PrintErr(err, "eoor occured in sending stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetPermissions(emp *emptypb.Empty, stream companypb.CompanyService_GetPermissionsServer) error {

	res, err := comp.Usecase.GetPermissions()
	if err != nil {
		helpers.PrintErr(err, "eroror occucres at GetPermissions usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.Permission{
			ID:         uint32(v.ID),
			Permission: v.Permission,
		}); err != nil {
			helpers.PrintErr(err, "eoor occured in sending stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) AttachRoleWithPermisssions(ctx context.Context, req *companypb.AttachRoleWithPermisssionsReq) (*emptypb.Empty, error) {

	if err := comp.Usecase.AttachRolewithPremission(entities.CompanyRoles{
		CompanyID:    req.CompanyID,
		RoleID:       uint(req.RoleID),
		PermissionID: uint(req.PermissionID),
	}); err != nil {
		helpers.PrintErr(err, "error occured at AttachRolewithPremission usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetAttachedRoleswithPermissions(req *companypb.GetAttachedRoleswithPermissionsReq, stream companypb.CompanyService_GetAttachedRoleswithPermissionsServer) error {

	res, err := comp.Usecase.GetRolesWithPermissions(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "errro occured at GetRolesWithPermissions usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetAttachedRoleswithPermissionsRes{
			ID:           uint32(v.ID),
			RoleID:       uint32(v.RoleID),
			PermissionID: uint32(v.PermissionID),
		}); err != nil {
			helpers.PrintErr(err, "error at sending stream")
			return err
		}
	}

	return nil
}
