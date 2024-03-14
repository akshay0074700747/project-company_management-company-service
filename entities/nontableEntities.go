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

type LogintoCompanyUsecaseTmp struct {
	CompanyID     string
	PermisssionID uint
	RoleID        uint
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

type GetPastProjectsUsecase struct {
	ClientID  string
	ProjectID string
}

type GetClientsUsecase struct {
	ClientID  string
	ProjectID []string
}

type GetRevenueGeneratedUsecase struct {
	ProjectID string
	ClientID  string
	Revenue   uint
}

type UpdateRevenueStatusUsecase struct {
	ProjectID  string
	ClientID   string
	IsRecieved bool
	CompanyID  string
}

type GetJobApplicationsofUserUsecase struct {
	ApplicatioID string
	Role         string
	IsVerified   bool
	CompanyID    string
}

type Responce struct {
	StatusCode int         `json:"StatusCode,omitempty"`
	Message    string      `json:"Message,omitempty"`
	Error      error       `json:"Error,omitempty"`
	Data       interface{} `json:"Data,omitempty"`
}

type UpdateAssetID struct {
	TransactionID string `json:"TransactionID"`
	UserID        string `json:"UserID"`
	AssetID       string `json:"AssetID"`
}
