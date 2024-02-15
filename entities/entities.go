package entities

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
}

type CompanyPhone struct {
	CompanyID string `gorm:"foreignKey:CompanyID;references:credentials(company_id)"`
	Phone     string `gorm:"unique;not null"`
}
