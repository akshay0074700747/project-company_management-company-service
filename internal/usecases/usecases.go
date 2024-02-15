package usecases

import (
	"errors"

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

func (comp *CompanyUseCases) RegisterCompany(req entities.CompanyResUsecase) (entities.CompanyResUsecase, error) {

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

	resCreds, err := comp.Adapter.InsertCompanyCredentials(req.CompCred)
	if err != nil {
		helpers.PrintErr(err, "error occured at InsertCompanyCredentials adapter")
	}

	var emails []entities.CompanyEmail
	var phones []entities.CompanyPhone

	for _, v := range req.Email {
		emails = append(emails, entities.CompanyEmail{
			Email:     v,
			CompanyID: resCreds.CompanyID,
		})
	}

	for _, v := range req.Phones {
		phones = append(phones, entities.CompanyPhone{
			Phone:     v,
			CompanyID: resCreds.CompanyID,
		})
	}

	_, err = comp.Adapter.InsertEmail(emails)
	if err != nil {
		helpers.PrintErr(err, "errror happened at InsertEmail adapter")
		return entities.CompanyResUsecase{}, err
	}

	_, err = comp.Adapter.InsertPhone(phones)
	if err != nil {
		helpers.PrintErr(err, "errror happened at InsertPhone adapter")
		return entities.CompanyResUsecase{}, err
	}

	req.Address.CompanyID = resCreds.CompanyID
	_, err = comp.Adapter.InsertAddress(req.Address)
	if err != nil {
		helpers.PrintErr(err, "errror happened at InsertAddress adapter")
		return entities.CompanyResUsecase{}, err
	}

	req.CompCred.CompanyID = resCreds.CompanyID

	return req, nil
}

func (comp *CompanyUseCases) AttachRolewithPremission(req entities.CompanyRoles) error {

	if err := comp.Adapter.AttachCompanyRoleAndPermissions(req); err != nil {
		helpers.PrintErr(err, "error occured at AttachCompanyRoleAndPermissions adapter")
		return err
	}

	return nil
}

func (comp *CompanyUseCases) AddMember(req entities.CompanyMembers) error {

	exists, err := comp.Adapter.IsMemberExists(req.MemberID)
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
