package usecases

import (
	"errors"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"github.com/akshay0074700747/project-company_management-company-service/helpers"
	"github.com/akshay0074700747/project-company_management-company-service/internal/adapters"
)

type CompanyUseCases struct {
	Adapter adapters.CompanyAdapterInterfaces
}

func NewCompanyUseCases(adapter adapters.CompanyAdapterInterfaces) *CompanyUseCases {
	return &CompanyUseCases{
		Adapter: adapter,
	}
}

func (comp *CompanyUseCases) RegisterCompany(req entities.CompanyResUsecase, ownerId string) (entities.CompanyResUsecase, error) {

	if req.CompCred.Name == "" {
		return entities.CompanyResUsecase{}, errors.New("the name cannot be empty")
	}

	if req.CompCred.CompanyUsername == "" {
		return entities.CompanyResUsecase{}, errors.New("the userNmae cannot be empty")
	}

	exists, err := comp.Adapter.IsCompanyUsernameExists(req.CompCred.CompanyUsername)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsCompanyUsernameExists adapter")
		return entities.CompanyResUsecase{}, err
	}

	if exists {
		return entities.CompanyResUsecase{}, errors.New("the username already exists")
	}

	req.CompCred.CompanyID = helpers.GenUuid()

	var resCreds entities.CompanyResUsecase

	resCreds.CompCred, resCreds.Email, resCreds.Phones, resCreds.Address, err = comp.Adapter.InsertCompanyCredentials(req.CompCred, req.Email, req.Phones, req.Address)
	if err != nil {
		helpers.PrintErr(err, "error occured at InsertCompanyCredentials adapter")
		return resCreds, err
	}

	return resCreds, nil
}

func (comp *CompanyUseCases) AttachRolewithPremission(req entities.CompanyRoles) error {

	if err := comp.Adapter.AttachCompanyRoleAndPermissions(req); err != nil {
		helpers.PrintErr(err, "error occured at AttachCompanyRoleAndPermissions adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AddMember(req entities.CompanyMembers) error {

	exists, err := comp.Adapter.IsMemberExists(req.MemberID, req.CompanyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsMemberExists adapter")
		return err
	}

	if exists {
		return errors.New("the user already exists")
	}

	exists, err = comp.Adapter.IsRoleIDExists(req.RoleID)
	if err != nil {
		helpers.PrintErr(err, "error occured at IsRoleIDExists adapter")
		return err
	}

	if !exists {
		return errors.New("the roleid doesnt exist")
	}

	if err = comp.Adapter.AddMember(req); err != nil {
		helpers.PrintErr(err, "error occured at AddMember adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetRolesWithPermissions(compID string) ([]entities.CompanyRoles, error) {

	res, err := comp.Adapter.GetRoleWithPermissionIDs(compID)
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetRoleWithPermissionIDs adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetPermissions() ([]entities.Permissions, error) {

	res, err := comp.Adapter.GetPermissions()
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetPermissions adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetCompanyTypes() ([]entities.CompanyTypes, error) {

	res, err := comp.Adapter.GetCompanyTypes()
	if err != nil {
		helpers.PrintErr(err, "eroor occured at GetCompanyTypes adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) AddCompanyType(req entities.CompanyTypes) error {

	if err := comp.Adapter.AddCompanyType(req); err != nil {
		helpers.PrintErr(err, "eroor occured at AddCompanyType adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AddPermissions(req entities.Permissions) error {

	if err := comp.Adapter.AddPermissions(req); err != nil {
		helpers.PrintErr(err, "eroor occured at AddPermissions adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetCompanyDetails(companyID string) (entities.ComapnyDetailsUsecase, error) {

	cred, err := comp.Adapter.GetCompanyDetails(companyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetProjectDetails adapter")
		return entities.ComapnyDetailsUsecase{}, err
	}
	mem, err := comp.Adapter.GetNoofMembers(companyID)
	if err != nil {
		helpers.PrintErr(err, "error at GetNoofMembers adapter")
		return entities.ComapnyDetailsUsecase{}, err
	}

	return entities.ComapnyDetailsUsecase{
		ComapanyID:      cred.CompanyID,
		CompanyUsername: cred.CompanyUsername,
		Name:            cred.Name,
		Members:         mem,
	}, nil
}

func (comp *CompanyUseCases) GetCompanyMembers(companyID string) ([]entities.GetCompanyEmployeesUsecase, error) {

	res, err := comp.Adapter.GetCompanyMembers(companyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetCompanyMembers usecase")
		return nil, err
	}

	return res, nil
}

func (company *CompanyUseCases) AddMemberStatueses(status string) error {

	if err := company.Adapter.AddMemberStatueses(status); err != nil {
		helpers.PrintErr(err, "error occured at AddMemberStatueses")
		return err
	}

	return nil
}

func (company *CompanyUseCases) GetAverageSalaryperRole(compID string) ([]entities.AverageSalaryperRoleUsecase, error) {

	res, err := company.Adapter.GetAverageSalaryperRole(compID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetAverageSalaryperRole adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) RaiseProblem(req entities.Problems) error {

	if err := comp.Adapter.RaiseProblem(req); err != nil {
		helpers.PrintErr(err, "error occured at RaiseProblem adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetProblems(companyID string) ([]entities.Problems, error) {

	problems, err := comp.Adapter.GetProblems(companyID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetProblems adapter")
		return nil, err
	}

	return problems, nil
}

func (comp *CompanyUseCases) InsertVisitors(companyUsername, visitorID string) error {

	companyID, err := comp.Adapter.GetCompanyIDFromName(companyUsername)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetCompanyIDFromName adapter")
		return err
	}

	timeVisited := time.Now()

	if err = comp.Adapter.InsertVisitors(entities.Visitors{
		CompanyID:   companyID,
		VisitorID:   visitorID,
		VisitedTime: timeVisited,
	}); err != nil {
		helpers.PrintErr(err, "error happened at InsertVisitors adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) GetVisitorsWithinTimeframe(compID string, from time.Time, to time.Time) ([]entities.Visitors, error) {

	res, err := comp.Adapter.GetVisitorsWithinTimeframe(compID, from, to)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetVisitorsWithinTimeframe adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetVisitors(compID string) ([]entities.Visitors, error) {

	res, err := comp.Adapter.GetVisitors(compID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetVisitors adapter")
		return nil, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetProfileViews(compID string, from time.Time, to time.Time) (int, error) {

	res, err := comp.Adapter.GetProfileViews(compID, from, to)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetProfileViews adapter")
		return 0, err
	}

	return res, nil
}

func (comp *CompanyUseCases) SalaryIncrementofEmployee(compID, userID string, increment int) error {

	if err := comp.Adapter.SalaryIncrementofEmployee(compID, userID, increment); err != nil {
		helpers.PrintErr(err, "error happened at SalaryIncrementofEmployee adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) SalaryIncrementofRole(compID string, roleId uint, increment int) error {

	if err := comp.Adapter.SalaryIncrementofRole(compID, roleId, increment); err != nil {
		helpers.PrintErr(err, "error happened at SalaryIncrementofRole adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) LogintoCompany(compUsername, userID string) (entities.LogintoCompanyUsecase, error) {

	compID, err := comp.Adapter.GetCompanyIDFromName(compUsername)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetCompanyIDFromName adapter")
		return entities.LogintoCompanyUsecase{}, err
	}

	res, err := comp.Adapter.LogintoCompany(compID, userID)
	if err != nil {
		helpers.PrintErr(err, "error happened at LogintoCompany adapter")
		return entities.LogintoCompanyUsecase{}, err
	}

	return res, nil
}

func (comp *CompanyUseCases) GetEmployeeLeaderBoard(companyID string) ([]entities.GetEmployeeLeaderBoardUsecase, error) {

	res, err := comp.Adapter.GetEmployeeLeaderBoard(companyID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetEmployeeLeaderBoard adapter")
		return nil, err
	}

	return res, nil
}
