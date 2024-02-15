package entities

type CompanyResUsecase struct {
	CompCred Credentials
	Email    []string
	Phones   []string
	Address  CompanyAddress
}
