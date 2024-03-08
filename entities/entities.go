package entities

import "time"

type Credentials struct {
	CompanyID       string `gorm:"primaryKey"`
	CompanyUsername string `gorm:"unique"`
	Name            string
	TypeID          uint `gorm:"foreignKey:TypeID;references:company_types(id)"`
	// Aim             string
}

type CompanyTypes struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"unique;not null"`
}

type CompanyAddress struct {
	CompanyID  string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	StreetNo   uint
	StreetName string
	PinNo      uint
	District   string
	State      string
	Nation     string
}

type CompanyEmail struct {
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Email     string `gorm:"unique;not null"`
}

type Permissions struct {
	ID         uint   `gorm:"primaryKey"`
	Permission string `gorm:"unique;not null"`
}

type CompanyRoles struct {
	ID           uint   `gorm:"primaryKey"`
	CompanyID    string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	RoleID       uint   `gorm:"not null"`
	PermissionID uint   `gorm:"foreignKey:PermissionID;references:permissions(id)"`
}

type CompanyMembers struct {
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	MemberID  string `gorm:"not null"`
	RoleID    uint   `gorm:"foreignKey:RoleID;references:company_roles(id)"`
	StatusID  uint   `gorm:"foreignKey:StatusID;references:member_statuses(id)"`
	Salary    int
}

type MemberStatus struct {
	ID     uint   `gorm:"primaryKey"`
	Status string `gorm:"unique"`
}

type CompanyPhone struct {
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Phone     string `gorm:"unique;not null"`
}

type Owners struct {
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	OwnerID   string
}

type Problems struct {
	ID                 uint   `gorm:"primaryKey"`
	CompanyID          string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Problem            string
	RaisedBy           string
	AssignedEmployeeID string
	IsResolved         bool
}

type Visitors struct {
	CompanyID   string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	VisitorID   string
	VisitedTime time.Time
}

type Clients struct {
	ID        uint `gorm:"primaryKey"`
	ClientID  string
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
}

type ClientsWithProjects struct {
	ClientID   uint `gorm:"foreignKey:ClientID;references:clients(id)"`
	ProjectID  string
	Contract   uint
	IsRecieved bool `gorm:"default:false"`
}

type CompanyPolicies struct {
	CompanyID          string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	MaxleavesPerMonth  uint32
	PayDay             uint
	WorkingHoursPerday uint32
}

type PayRoll struct {
	CompanyID     string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	EmployeeID    string
	IsPayed       bool
	TransactionID string
	Date          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Leaves struct {
	ID          uint `gorm:"primaryKey"`
	EmployeeID  string
	CompanyID   string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Description string
	Date        time.Time
	IsAllowed   bool `gorm:"default:false"`
}

type Jobs struct {
	JobID               string `gorm:"primaryKey"`
	CompanyID           string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Role                string
	Vacancy             uint32
	Description         string
	MinExperiance       uint32
	MinExpectedCTC      float32
	MaxExpectedCTC      float32
	IsRemote            bool
	AddressID           uint `gorm:"foreignKey:AddressID;references:addresses(id)"`
	TotalPersonsApplied uint `gorm:"-"`
}

type Address struct {
	ID         uint   `gorm:"primaryKey" json:"ID"`
	StreetNo   uint   `json:"StreetNo"`
	StreetName string `json:"StreetName"`
	PinNo      uint   `json:"PinNo"`
	District   string `json:"District"`
	State      string `json:"State"`
	Nation     string `json:"Nation"`
}

type JobApplications struct {
	ApplicationID      string  `gorm:"primaryKey" json:"ApplicationID"`
	UserID             string  `json:"UserID"`
	JobID              string  `gorm:"foreignKey:JobID;references:jobs(job_id)" json:"JobID"`
	Name               string  `json:"Name"`
	Email              string  `json:"Email"`
	Mobile             string  `json:"Mobile"`
	AddressofApplicant Address `gorm:"-" json:"AddressofApplicant"`
	HighestEducation   string  `json:"HighestEducation"`
	Nationality        string  `json:"Nationality"`
	Experiance         uint32  `json:"Experiance"`
	CurrentCTC         float32 `json:"CurrentCTC"`
	Resume             []byte  `gorm:"-" json:"Resume"`
	ResumeID           string
	FileName           string `gorm:"-" json:"FileName"`
	AddressID          uint
	IsShortlisted      bool `gorm:"default:false"`
}

type ScheduledInterviews struct {
	ID            uint   `gorm:"primaryKey"`
	ApplicationID string `gorm:"foreignKey:ApplicationID;references:job_applications(application_id)" json:"JobID"`
	Date          time.Time
	Description   string
	Time          string
}
