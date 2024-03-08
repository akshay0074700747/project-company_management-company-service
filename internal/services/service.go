package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/akshay0074700747/project-company_management-company-service/helpers"
	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/companypb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/projectpb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CompanyServiceServer struct {
	UserConn    userpb.UserServiceClient
	ProjectConn projectpb.ProjectServiceClient
	Usecase     usecases.CompanyUsecaseInterfaces
	companypb.UnimplementedCompanyServiceServer
}

func NewProjectServiceServer(usecase usecases.CompanyUsecaseInterfaces, addr, projectAddr string) *CompanyServiceServer {
	userRes, _ := helpers.DialGrpc(addr)
	projectRes, _ := helpers.DialGrpc(projectAddr)
	return &CompanyServiceServer{
		Usecase:     usecase,
		UserConn:    userpb.NewUserServiceClient(userRes),
		ProjectConn: projectpb.NewProjectServiceClient(projectRes),
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
		Salary:    int(req.CTC),
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

	res, err := comp.Usecase.GetCompanyMembers(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetCompanyMembers usecase")
		return err
	}

	streeam, err := comp.UserConn.GetStreamofUserDetails(context.TODO())
	if err != nil {
		helpers.PrintErr(err, "error at stream GetStreamofUserDetails")
		return err
	}

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
			UserId:     details.UserID,
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

	isOwnerbool, err := company.Usecase.IsOwner(res.CompanyID, req.UserID)
	if err != nil {
		helpers.PrintErr(err, "error at IsOwner usecase")
		return nil, err
	}

	fmt.Println(isOwnerbool)

	if isOwnerbool {
		return &companypb.LogintoCompanyRes{
			CompanyID:  res.CompanyID,
			Permission: "ROOT",
			Role:       "OWNER",
		}, nil
	}

	if res.CompanyID == "" || res.Permisssion == "" || res.RoleID == 0 {
		return nil, errors.New("the compnayusername or memberID is not valid")
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

	var isAssigned bool
	for _, v := range problems {

		if v.AssignedEmployeeID != "" {
			isAssigned = true
		} else {
			isAssigned = false
		}
		if err = stream.Send(&companypb.GetProblemsRes{
			ProblemID:  uint32(v.ID),
			Problem:    v.Problem,
			RaisedBy:   v.RaisedBy,
			IsResolved: v.IsResolved,
			IsAssigned: isAssigned,
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

func (company *CompanyServiceServer) GetPermission(ctx context.Context, req *companypb.GetPermisssionReq) (*companypb.GetPermisssionRes, error) {

	permission, err := company.Usecase.GetPermission(uint(req.ID))
	if err != nil {
		helpers.PrintErr(err, "error happened at GetPermission usecase")
		return nil, err
	}

	return &companypb.GetPermisssionRes{
		Permission: permission,
	}, nil
}

func (comp *CompanyServiceServer) IsEmployeeExists(ctx context.Context, req *companypb.IsEmployeeExistsReq) (*companypb.IsEmployeeExistsRes, error) {

	exists, err := comp.Usecase.IsEmployeeExists(req.EmployeeID, req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at IsEmployeeExists usecase")
		return nil, err
	}

	return &companypb.IsEmployeeExistsRes{Exists: exists}, nil
}

func (company *CompanyServiceServer) AddClient(ctx context.Context, req *companypb.AddClientReq) (*emptypb.Empty, error) {

	details, err := company.UserConn.GetByEmail(ctx, &userpb.GetByEmailReq{
		Email: req.Email,
	})
	if err != nil {
		helpers.PrintErr(err, "error happened at GetByEmail")
		return nil, err
	}

	if err := company.Usecase.InsertIntoClients(entities.Clients{
		ClientID:  details.UserID,
		CompanyID: req.CompanyID,
	}); err != nil {
		helpers.PrintErr(err, "error happened at InsertIntoClients usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) AssociateClientWithProject(ctx context.Context, req *companypb.AssociateClientWithProjectReq) (*emptypb.Empty, error) {

	if err := company.Usecase.AttachClientwithProject(entities.Clients{
		ClientID:  req.ClientID,
		CompanyID: req.CompanyID,
	}, req.ProjectID, uint(req.Contract)); err != nil {
		helpers.PrintErr(err, "error happened at AttachClientwithProject usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) GetPastProjects(req *companypb.GetProjectsReq, stream companypb.CompanyService_GetPastProjectsServer) error {

	res, err := company.Usecase.GetPastProjects(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetPastProjects usecase")
		return err
	}

	streaam, err := company.ProjectConn.GetStreamofProjectDetails(context.TODO())

	for _, v := range res {
		if err := streaam.Send(&projectpb.GetStreamofProjectDetailsReq{
			ProjectID: v.ProjectID,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}

		details, err := streaam.Recv()
		if err != nil {
			helpers.PrintErr(err, "error happened at recieving from stream")
			return err
		}

		if err := stream.Send(&companypb.GetProjectsRes{
			ProjectID:   details.ProjectID,
			Description: details.Aim,
			ProjectName: details.ProjectUsername,
			ClientID:    v.ClientID,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (companay *CompanyServiceServer) GetClients(req *companypb.GetClientsReq, stream companypb.CompanyService_GetClientsServer) error {

	res, err := companay.Usecase.GetClients(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetClients usecase")
		return err
	}

	streaam, err := companay.UserConn.GetStreamofUserDetails(context.TODO())
	if err != nil {
		helpers.PrintErr(err, "error happened at GetStreamofUserDetails")
		return err
	}
	for _, v := range res {
		if err := streaam.Send(&userpb.GetUserDetailsReq{
			UserID: v.ClientID,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}

		details, err := streaam.Recv()
		if err != nil {
			helpers.PrintErr(err, "error happened at recieving from stream")
			return err
		}

		if err = stream.Send(&companypb.GetClientsRes{
			ClientID:   v.ClientID,
			ProjectIDs: v.ProjectID,
			Name:       details.Name,
			Email:      details.Email,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) GetRevenueGenerated(req *companypb.GetRevenueGeneratedReq, stream companypb.CompanyService_GetRevenueGeneratedServer) error {

	res, err := company.Usecase.GetRevenuesGenerated(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happenened at GetRevenuesGenerated usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetRevenueGeneratedRes{
			ProjectID: v.ProjectID,
			ClientID:  v.ClientID,
			Revenue:   uint32(v.Revenue),
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) UpdateRevenueStatus(ctx context.Context, req *companypb.UpdateRevenueStatusReq) (*emptypb.Empty, error) {

	if err := company.Usecase.UpdateRevenueStatus(entities.UpdateRevenueStatusUsecase{
		ProjectID:  req.ProjectID,
		ClientID:   req.ClientID,
		IsRecieved: req.IsRecieved,
		CompanyID:  req.CompanyID,
	}); err != nil {
		helpers.PrintErr(err, "error happened at UpdateRevenueStatus usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) AttachCompanyPolicies(ctx context.Context, req *companypb.AttachCompanyPoliciesReq) (*emptypb.Empty, error) {

	if err := company.Usecase.UpdateCompanyPolicies(entities.CompanyPolicies{
		CompanyID:          req.CompanyID,
		MaxleavesPerMonth:  req.MaxleavesPerMonth,
		PayDay:             uint(req.PayDay),
		WorkingHoursPerday: req.WorkingHoursPerday,
	}); err != nil {
		helpers.PrintErr(err, "eroror happeneded at UpdateCompanyPolicies usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) UpdatePaymentStatusofEmployee(ctx context.Context, req *companypb.UpdatePaymentStatusofEmployeeReq) (*emptypb.Empty, error) {

	if err := company.Usecase.UpdatePayRollofEmployee(entities.PayRoll{
		CompanyID:     req.CompanyID,
		EmployeeID:    req.EmployeeID,
		TransactionID: req.TransactionID,
		IsPayed:       req.IsPayed,
	}); err != nil {
		helpers.PrintErr(err, "eroror happened at UpdatePayRollofEmployee usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) AssignProblem(ctx context.Context, req *companypb.AssignProblemReq) (*emptypb.Empty, error) {

	if err := company.Usecase.AssignProblemToEmployee(req.EmployeeID, uint(req.ProblemID)); err != nil {
		helpers.PrintErr(err, "eroro happened at AssignProblemToEmployee usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) ResolveProblem(ctx context.Context, req *companypb.ResolveProblemReq) (*emptypb.Empty, error) {

	if err := company.Usecase.ResolveProblem(uint(req.ProblemID), req.ResolverID); err != nil {
		helpers.PrintErr(err, "eroro happened at ResolveProblem usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) ApplyForLeave(ctx context.Context, req *companypb.ApplyForLeaveReq) (*emptypb.Empty, error) {

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helpers.PrintErr(err, "error ahppenede at parsing to time")
		return nil, err
	}

	if err = company.Usecase.ApplyforLeave(entities.Leaves{
		EmployeeID:  req.EmployeeID,
		CompanyID:   req.CompanyID,
		Date:        date,
		Description: req.Description,
	}); err != nil {
		helpers.PrintErr(err, "eroror happenede at ApplyforLeave usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (company *CompanyServiceServer) GetEmployeeLeaveRequests(req *companypb.GetEmployeeLeaveRequestsReq, stream companypb.CompanyService_GetEmployeeLeaveRequestsServer) error {

	res, err := company.Usecase.GetAppliedLeaves(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetAppliedLeaves usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetEmployeeLeaveRequestsRes{
			LeaveID:     uint32(v.ID),
			Date:        v.Date.String(),
			EmployeeID:  v.EmployeeID,
			Description: v.Description,
		}); err != nil {
			helpers.PrintErr(err, "eroror happened at sending to stream")
			return err
		}
	}

	return nil
}

func (company *CompanyServiceServer) DecideEmployeeLeave(ctx context.Context, req *companypb.DecideEmployeeLeaveRequest) (*emptypb.Empty, error) {

	if err := company.Usecase.GrantLeave(uint(req.LeaveID), req.IsAllowed); err != nil {
		helpers.PrintErr(err, "error happened at GrantLeave usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetLeaves(req *companypb.GetLeavesReq, stream companypb.CompanyService_GetLeavesServer) error {

	res, err := comp.Usecase.GetLeaves(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happenede at GetLeaves usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetLeavesRes{
			EmployeeID:  v.EmployeeID,
			Date:        v.Date.String(),
			Description: v.Description,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetStreamofClients(stream companypb.CompanyService_GetStreamofClientsServer) error {

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			helpers.PrintErr(err, "eroro occured at recieving from stream")
			return err
		}

		id, err := comp.Usecase.GetClientID(req.ProjectID)
		if err != nil {
			helpers.PrintErr(err, "error happened at GetClientID usecase")
			return err
		}

		if err = stream.Send(&companypb.GetStreamofClientsRes{
			ClientID: id,
		}); err != nil {
			helpers.PrintErr(err, "eroro happened at sending to sttream")
			return err
		}
	}

	return nil
}

func (project *CompanyServiceServer) PostJobs(ctx context.Context, req *companypb.PostJobsReq) (*emptypb.Empty, error) {

	fmt.Println(req)

	if err := project.Usecase.PostJob(entities.Address{
		StreetNo:   uint(req.JobLocation.StreetNo),
		StreetName: req.JobLocation.Street,
		PinNo:      uint(req.JobLocation.PinNo),
		District:   req.JobLocation.District,
		State:      req.JobLocation.State,
		Nation:     req.JobLocation.Nation,
	}, entities.Jobs{
		CompanyID:      req.CompanyID,
		Role:           req.Role,
		Vacancy:        req.Vacancy,
		Description:    req.Description,
		MinExperiance:  req.MinExperiance,
		MinExpectedCTC: req.MinExpectedCTC,
		MaxExpectedCTC: req.MaxExpectedCTC,
		IsRemote:       req.IsRemote,
	}); err != nil {
		helpers.PrintErr(err, "errror happened at PostJobs usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetJobsofCompany(req *companypb.GetJobsofCompanyReq, stream companypb.CompanyService_GetJobsofCompanyServer) error {

	res, err := comp.Usecase.GetJobsofCompany(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetJobsofCompany usecase")
		return err
	}

	for _, v := range res {
		if err := stream.Send(&companypb.GetJobsofCompanyRes{
			JobID:               v.JobID,
			Role:                v.Role,
			Description:         v.Description,
			Vacancy:             v.Vacancy,
			MinExperiance:       v.MinExperiance,
			MinExpectedCTC:      v.MinExpectedCTC,
			MaxExpectedCTC:      v.MaxExpectedCTC,
			IsRemote:            v.IsRemote,
			TotalPersonsApplied: uint32(v.TotalPersonsApplied),
		}); err != nil {
			helpers.PrintErr(err, "eroror happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetJobApplications(req *companypb.GetJobApplicationsReq, stream companypb.CompanyService_GetJobApplicationsServer) error {

	res, err := comp.Usecase.GetApplicationsforJob(req.JobID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetApplicationsforJob usecase")
		return err
	}

	for _, v := range res {
		if err := stream.Send(&companypb.GetJobApplicationsRes{
			ApplicationID:    v.ApplicationID,
			Name:             v.Name,
			Email:            v.Email,
			Mobile:           v.Mobile,
			ResumeID:         v.ResumeID,
			HighestEducation: v.HighestEducation,
			Nationality:      v.Nationality,
			Experiance:       v.Experiance,
			CurrentCTC:       v.CurrentCTC,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) ShortlistApplications(ctx context.Context, req *companypb.ShortlistApplicationsReq) (*emptypb.Empty, error) {

	if err := comp.Usecase.ShortlistApplications(req.ApplicationID); err != nil {
		helpers.PrintErr(err, "error happeed at ShortlistApplications usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) ScheduleInterview(ctx context.Context, req *companypb.ScheduleInterviewReq) (*emptypb.Empty, error) {

	timme, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helpers.PrintErr(err, "error happened at parsing")
		return nil, err
	}

	if err := comp.Usecase.ScheduleInterviews(entities.ScheduledInterviews{
		ApplicationID: req.ApplicationID,
		Date:          timme,
		Description:   req.Description,
		Time:          req.Time,
	}); err != nil {
		helpers.PrintErr(err, "error happened at ScheduleInterviews usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetScheduledInterviews(req *companypb.GetScheduledInterviewsReq, stream companypb.CompanyService_GetScheduledInterviewsServer) error {

	res, err := comp.Usecase.GetScheduledInterviews(req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetScheduledInterviews usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetScheduledInterviewsRes{
			ApplicationID: v.ApplicationID,
			Date:          v.Date.String(),
			Description:   v.Description,
			Time:          v.Time,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetDetailsofApplicationByID(ctx context.Context, req *companypb.GetDetailsofApplicationByIDReq) (*companypb.GetDetailsofApplicationByIDRes, error) {

	res, err := comp.Usecase.GetDetialsodApplicationbyID(req.ApplicationID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetDetialsodApplicationbyID usecase")
		return nil, err
	}

	return &companypb.GetDetailsofApplicationByIDRes{
		ApplicationID:    res.ApplicationID,
		Name:             res.Name,
		Email:            res.Email,
		Mobile:           res.Mobile,
		ResumeID:         res.ResumeID,
		HighestEducation: res.HighestEducation,
		Nationality:      res.Nationality,
		Experiance:       res.Experiance,
		CurrentCTC:       res.CurrentCTC,
	}, nil
}

func (comp *CompanyServiceServer) GetScheduledInterviewsofUser(req *companypb.GetScheduledInterviewsofUserReq, stream companypb.CompanyService_GetScheduledInterviewsofUserServer) error {

	res, err := comp.Usecase.GetScheduledInterviewsofUser(req.UserID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetScheduledInterviewsofUser usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetScheduledInterviewsofUserRes{
			ApplicationID: v.ApplicationID,
			Date:          v.Date.String(),
			Description:   v.Description,
			Time:          v.Time,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) RescheduleInterview(ctx context.Context, req *companypb.RescheduleInterviewReq) (*emptypb.Empty, error) {

	dte, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helpers.PrintErr(err, "erorr happened at parsing")
		return nil, err
	}

	if err := comp.Usecase.RescheduleInterview(entities.ScheduledInterviews{
		ApplicationID: req.ApplicationID,
		Date:          dte,
		Description:   req.Description,
		Time:          req.Time,
	}); err != nil {
		helpers.PrintErr(err, "error happened at RescheduleInterview usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (comp *CompanyServiceServer) GetShortlistedApplications(req *companypb.GetShortlistedApplicationsReq, stream companypb.CompanyService_GetShortlistedApplicationsServer) error {

	res, err := comp.Usecase.GetShortlistedApplications(req.JobID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetShortlistedApplications usecase")
		return err
	}

	for _, v := range res {

		if err := stream.Send(&companypb.GetShortlistedApplicationsRes{
			ApplicationID:    v.ApplicationID,
			Name:             v.Name,
			Email:            v.Email,
			Mobile:           v.Mobile,
			ResumeID:         v.ResumeID,
			HighestEducation: v.HighestEducation,
			Nationality:      v.Nationality,
			Experiance:       v.Experiance,
			CurrentCTC:       v.CurrentCTC,
		}); err != nil {
			helpers.PrintErr(err, "errror happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetJobs(req *companypb.GetJobsReq, stream companypb.CompanyService_GetJobsServer) error {

	jobs, err := comp.Usecase.GetJobs(req.CompanyID, req.Role)
	if err != nil {
		helpers.PrintErr(err, "eroror happened at GetJobs usecase")
		return err
	}

	for _, v := range jobs {

		if err = stream.Send(&companypb.GetJobsRes{
			JobID:          v.JobID,
			CompanyID:      v.CompanyID,
			Role:           v.Role,
			Vacancy:        v.Vacancy,
			Description:    v.Description,
			MinExperiance:  v.MinExperiance,
			MinExpectedCTC: v.MinExpectedCTC,
			MaxExpectedCTC: v.MaxExpectedCTC,
			IsRemote:       v.IsRemote,
		}); err != nil {
			helpers.PrintErr(err, "error happend at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetAllJobApplicationsofUser(req *companypb.GetAllJobApplicationsofUserReq, stream companypb.CompanyService_GetAllJobApplicationsofUserServer) error {

	res, err := comp.Usecase.GetJobApplicationsofUser(req.UserID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetJobApplicationsofUser usecase")
		return err
	}

	var status string
	for _, v := range res {

		if v.IsVerified {
			status = "ShortListed"
		} else {
			status = "Not ShortListed"
		}

		if err := stream.Send(&companypb.GetAllJobApplicationsofUserRes{
			ApplicationID: v.ApplicatioID,
			CompanyID:     v.CompanyID,
			Role:          v.Role,
			Status:        status,
		}); err != nil {
			helpers.PrintErr(err, "error happened at sending to stream")
			return err
		}
	}

	return nil
}

func (comp *CompanyServiceServer) GetAssignedProblems(req *companypb.GetAssignedProblemsReq, stream companypb.CompanyService_GetAssignedProblemsServer) error {

	res, err := comp.Usecase.GetAssignedProblems(req.CompanyID, req.UserID)
	if err != nil {
		helpers.PrintErr(err, "erorr happened at GetAssignedProblems usecase")
		return err
	}

	for _, v := range res {
		if err = stream.Send(&companypb.GetAssignedProblemsRes{
			ProblemID: uint32(v.ID),
			Problem:   v.Problem,
			RaisedBy:  v.RaisedBy,
		}); err != nil {
			helpers.PrintErr(err, "errror happened at sending to stream")
			return err
		}
	}

	return nil
}
