package entities

type CompanyResUsecase struct {
	CompCred Credentials
	Email    []string
	Phones   []string
	Address  CompanyAddress
}

type ComapnyDetailsUsecase struct {
	ComapanyID      string
	CompanyUsername string
	Name            string
	Members         uint
}

type AverageSalaryperRoleUsecase struct {
	RoleID uint
	Salary float32
}

type LogintoCompanyUsecase struct {
	CompanyID   string
	Permisssion string
	RoleID      uint
}

type ListCompaniesUsecase struct {
	CompanyID   string
	Aim         string
	Description string
	Employees   uint32
}

type GetEmployeeLeaderBoardUsecase struct {
	EmployeeID string
	Salary     int
	RoleID     uint
}

type GetCompanyEmployeesUsecase struct {
	UserID     string
	Permission string
	RoleID     uint
}
