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
	"google.golang.org/protobuf/types/known/timestamppb"
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
	}, req.OwnerID)
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
		return nil, err
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

func (comp *CompanyServiceServer) AddCompanyTypes(ctx context.Context, req *companypb.AddCompanyTypeReq) (*emptypb.Empty, error) {

	if err := comp.Usecase.AddCompanyType(entities.CompanyTypes{
		Type: req.Name,
	}); err != nil {
		helpers.PrintErr(err, "error at AddCompanyType usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) Permissions(ctx context.Context, req *companypb.AddPermissionReq) (*emptypb.Empty, error) {

	if err := comp.Usecase.AddPermissions(entities.Permissions{
		Permission: req.Name,
	}); err != nil {
		helpers.PrintErr(err, "error at AddPermissions usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetCompanyDetails(ctx context.Context, req *companypb.GetCompanyReq) (*companypb.GetCompanyDetailsRes, error) {

	if req.CompanyUsername != "" {
		if err := comp.Usecase.InsertVisitors(req.CompanyUsername, req.VisitorID); err != nil {
			helpers.PrintErr(err, "error at InsertVisitors usecase")
		}
	}
	res, err := comp.Usecase.GetCompanyDetails(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetCompanyDetails usecase")
		return nil, errors.New("cannot connect to the company service now")
	}

	return &companypb.GetCompanyDetailsRes{
		CompanyID:       res.ComapanyID,
		CompanyUsername: res.CompanyUsername,
		Members:         uint32(res.Members),
	}, nil
}

func (comp *CompanyServiceServer) GetCompanyEmployees(req *companypb.GetCompanyReq, stream companypb.CompanyService_GetCompanyEmployeesServer) error {

	res,err := comp.Usecase.GetCompanyMembers(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetCompanyMembers usecase")
		return err
	}

	streeam, err := comp.UserConn.GetStreamofUserDetails(context.TODO())

	for i := range res {

		if err = streeam.Send(&userpb.GetUserDetailsReq{
			UserID: res[i].UserID,
			RoleID: uint32(res[i].RoleID),
		}); err != nil {
			helpers.PrintErr(err, "error at sending to stream")
			return err
		}

		details, err := streeam.Recv()
		if err != nil {
			helpers.PrintErr(err, "error at recieving from stream")
			return err
		}

		if err := stream.Send(&companypb.GetCompanyEmployeesRes{
			Email:      details.Email,
			Name:       details.Name,
			Permission: res[i].Permission,
			Role:       details.Role,
		}); err != nil {
			helpers.PrintErr(err, "error at sending stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) LogintoCompany(ctx context.Context, req *companypb.LogintoCompanyReq) (*companypb.LogintoCompanyRes, error) {

	res, err := company.Usecase.LogintoCompany(req.CompanyUsername, req.UserID)
	if err != nil {
		helpers.PrintErr(err, "error at LogintoCompany usecase")
		return nil, err
	}
	role, err := company.UserConn.GetRole(ctx, &userpb.GetRoleReq{
		ID: uint32(res.RoleID),
	})
	if err != nil {
		helpers.PrintErr(err, "error at communicating with the user service")
		return nil, err
	}

	return &companypb.LogintoCompanyRes{
		CompanyID:  res.CompanyID,
		Permission: res.Permisssion,
		Role:       role.Role,
	}, nil
}

func (company *CompanyServiceServer) AddMemberStatus(ctx context.Context, req *companypb.MemberStatusReq) (*emptypb.Empty, error) {

	if err := company.Usecase.AddMemberStatueses(req.Status); err != nil {
		helpers.PrintErr(err, "error happeded att AddMemberStatueses")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) GetAverageSalaryperRole(req *companypb.GetAverageSalaryperRoleReq, stream companypb.CompanyService_GetAverageSalaryperRoleServer) error {

	res, err := company.Usecase.GetAverageSalaryperRole(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happeded att GetAverageSalaryperRole usecase")
		return err
	}

	streaam, err := company.UserConn.GetStreamofRoles(context.TODO())
	if err != nil {
		helpers.PrintErr(err, "error happeded att creating stream")
		return err
	}

	for _, v := range res {
		if err := streaam.Send(&userpb.GetStreamofRolesReq{
			RoleID: uint32(v.RoleID),
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending stream")
			return err
		}
	}
	resMap, err := streaam.CloseAndRecv()
	if err != nil {
		helpers.PrintErr(err, "error happeded att creating stream")
		return err
	}

	for _, v := range res {
		if err := stream.Send(&companypb.GetAverageSalaryperRoleRes{
			Role:   resMap.RoleIDsWithNames[uint32(v.RoleID)],
			Salary: uint32(v.Salary),
		}); err != nil {
			helpers.PrintErr(err, "error happeded att cresendingating stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) GetProblems(req *companypb.GetProblemsReq, stream companypb.CompanyService_GetProblemsServer) error {

	problems, err := company.Usecase.GetProblems(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happeded att GetProblems usecase")
		return err
	}

	for _, v := range problems {
		if err = stream.Send(&companypb.GetProblemsRes{
			Problem:  v.Problem,
			RaisedBy: v.RaisedBy,
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) GetProfileViews(ctx context.Context, req *companypb.GetProfileViewsReq) (*companypb.GetProfileViewsRes, error) {

	var res []entities.Visitors
	var err error
	if req.From != nil && req.To != nil {
		res, err = company.Usecase.GetVisitorsWithinTimeframe(req.CompanyID, req.From.AsTime(), req.To.AsTime())
		if err != nil {
			helpers.PrintErr(err, "error happeded att GetVisitorsWithinTimeframe usecase")
			return &companypb.GetProfileViewsRes{}, err
		}
	} else {
		res, err = company.Usecase.GetVisitors(req.CompanyID)
		if err != nil {
			helpers.PrintErr(err, "error happeded att GetVisitors usecase")
			return &companypb.GetProfileViewsRes{}, err
		}
	}

	return &companypb.GetProfileViewsRes{
		CompanyID: req.CompanyID,
		Views:     uint32(len(res)),
	}, nil
}

func (company *CompanyServiceServer) GetSalaryLeaderboard(req *companypb.GetSalaryLeaderboardReq, stream companypb.CompanyService_GetSalaryLeaderboardServer) error {

	res, err := company.Usecase.GetEmployeeLeaderBoard(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happeded att GetEmployeeLeaderBoard usecase")
		return err
	}

	streeam, err := company.UserConn.GetStreamofUserDetails(context.TODO())
	if err != nil {
		helpers.PrintErr(err, "error happeded att GetStreamofUserDetails usecase")
		return err
	}

	for _, v := range res {

		if err := streeam.Send(&userpb.GetUserDetailsReq{
			UserID: v.EmployeeID,
			RoleID: uint32(v.RoleID),
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending stream")
			return err
		}

		detaisl, err := streeam.Recv()
		if err != nil {
			helpers.PrintErr(err, "error happeded att recieving stream")
			return err
		}

		if err := stream.Send(&companypb.GetSalaryLeaderboardRes{
			EmployeeID: detaisl.UserID,
			Email:      detaisl.Email,
			Name:       detaisl.Name,
			Role:       detaisl.Role,
			Salary:     uint32(v.Salary),
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) RaiseProblem(ctx context.Context, req *companypb.RaiseProblemReq) (*emptypb.Empty, error) {

	if err := company.Usecase.RaiseProblem(entities.Problems{
		CompanyID: req.CompanyID,
		Problem:   req.Problem,
		RaisedBy:  req.UserID,
	}); err != nil {
		helpers.PrintErr(err, "error happeded att RaiseProblem usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) SalaryIncrementofEmployee(ctx context.Context, req *companypb.SalaryIncrementofEmployeeReq) (*emptypb.Empty, error) {

	if err := company.Usecase.SalaryIncrementofEmployee(req.CompanyID, req.EmployeeID, int(req.Increment)); err != nil {
		helpers.PrintErr(err, "error happeded att SalaryIncrementofEmployee usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) SalaryIncrementofRole(ctx context.Context, req *companypb.SalaryIncrementofRoleReq) (*emptypb.Empty, error) {

	if err := company.Usecase.SalaryIncrementofRole(req.CompanyID, uint(req.RoleID), int(req.Increment)); err != nil {
		helpers.PrintErr(err, "error happeded att SalaryIncrementofRole usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) GetVisitors(req *companypb.GetVisitorsReq, stream companypb.CompanyService_GetVisitorsServer) error {

	var res []entities.Visitors
	var err error
	if req.From != nil && req.To != nil {

		res, err = company.Usecase.GetVisitorsWithinTimeframe(req.CompanyID, req.From.AsTime(), req.To.AsTime())
		if err != nil {
			helpers.PrintErr(err, "error happeded att GetVisitorsWithinTimeframe usecase")
			return err
		}
	} else {
		res, err = company.Usecase.GetVisitors(req.CompanyID)
		if err != nil {
			helpers.PrintErr(err, "error happeded att GetVisitors usecase")
			return err
		}
	}

	streaam, err := company.UserConn.GetStreamofUserDetails(context.TODO())
	if err != nil {
		helpers.PrintErr(err, "error happeded att getting stream")
		return err
	}

	for _, v := range res {
		if err := streaam.Send(&userpb.GetUserDetailsReq{
			UserID: v.VisitorID,
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending to stream")
			return err
		}
		details, err := streaam.Recv()
		if err != nil {
			helpers.PrintErr(err, "error happeded att recieving stream")
			return err
		}
		if stream.Send(&companypb.GetVisitorsRes{
			Name:        details.Name,
			Email:       details.Email,
			VisitedTime: timestamppb.New(v.VisitedTime),
		}); err != nil {
			helpers.PrintErr(err, "error happeded att sending stream")
			return err
		}
	}

	return nil
}
