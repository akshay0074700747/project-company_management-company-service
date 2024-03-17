package test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/entities"
	mock_adapters "github.com/akshay0074700747/project-company_management-company-service/internal/adapters/mockAdapters"
	"github.com/akshay0074700747/project-company_management-company-service/internal/usecases"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCompany(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name                         string
		mockIscompanyusernameTaken   func(companyUsername string) (bool, error)
		mockInsertcompanyCredentials func(entities.Credentials, []string, []string, entities.CompanyAddress, string) (entities.Credentials, []string, []string, entities.CompanyAddress, error)
		credentialsRequest           entities.CompanyResUsecase
		ownerIDrequest               string
		wantError                    bool
		wantResult                   entities.CompanyResUsecase
	}{
		{
			name: "Success",
			mockIscompanyusernameTaken: func(companyUsername string) (bool, error) {
				return "itsunique" == companyUsername, nil
			},
			mockInsertcompanyCredentials: func(c entities.Credentials, s1, s2 []string, ca entities.CompanyAddress, s3 string) (entities.Credentials, []string, []string, entities.CompanyAddress, error) {
				return c, s1, s2, ca, nil
			},
			credentialsRequest: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "itsnotUnique",
					Name:            "NewNmae",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
			ownerIDrequest: "73ghdygf6734vud7g223",
			wantError:      false,
			wantResult: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "itsnotUnique",
					Name:            "NewNmae",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
		},
		{
			name: "Fail",
			mockIscompanyusernameTaken: func(companyUsername string) (bool, error) {
				return "itsunique" == companyUsername, nil
			},
			mockInsertcompanyCredentials: func(c entities.Credentials, s1, s2 []string, ca entities.CompanyAddress, s3 string) (entities.Credentials, []string, []string, entities.CompanyAddress, error) {
				return c, s1, s2, ca, nil
			},
			credentialsRequest: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "",
					Name:            "NewNmae",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
			ownerIDrequest: "73ghdygf6734vud7g223",
			wantError:      true,
			wantResult: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "",
					Name:            "NewNmae",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
		},
		{
			name: "Fail",
			mockIscompanyusernameTaken: func(companyUsername string) (bool, error) {
				return "itsunique" == companyUsername, nil
			},
			mockInsertcompanyCredentials: func(c entities.Credentials, s1, s2 []string, ca entities.CompanyAddress, s3 string) (entities.Credentials, []string, []string, entities.CompanyAddress, error) {
				return c, s1, s2, ca, nil
			},
			credentialsRequest: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "itsnotUnique",
					Name:            "",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
			ownerIDrequest: "73ghdygf6734vud7g223",
			wantError:      true,
			wantResult: entities.CompanyResUsecase{
				CompCred: entities.Credentials{
					CompanyUsername: "itsnotUnique",
					Name:            "",
					TypeID:          1,
				},
				Email:  []string{"sadhcksdj@gmail.com", "gsjj@gmail.com"},
				Phones: []string{"645271312", "8723948873409"},
				Address: entities.CompanyAddress{
					StreetNo:   23,
					StreetName: "hadsvkhjs",
					PinNo:      8439287,
					District:   "jhdksjks",
					State:      "cgdskhjcask",
					Nation:     "dskakhjlk",
				},
			},
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().IsCompanyUsernameExists(gomock.Any()).DoAndReturn(test.mockIscompanyusernameTaken).AnyTimes().Times(1)
			adapter.EXPECT().InsertCompanyCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(test.mockInsertcompanyCredentials).AnyTimes().Times(1)
		}
		regUsecase := usecases.NewCompanyUseCases(adapter)

		res, err := regUsecase.RegisterCompany(test.credentialsRequest, test.ownerIDrequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, res)
			res.CompCred.CompanyID = ""
			assert.Equal(t, test.wantResult, res)
		}
	}
}

func TestAddMember(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name               string
		mockIsMemberExists func(string, string) (bool, error)
		mockIsRoleIDExists func(uint) (bool, error)
		mockAddMember      func(entities.CompanyMembers) error
		memberRequest      entities.CompanyMembers
		isMemberExist      bool
		isRoleExist        bool
	}{
		{
			name: "Success",
			mockIsMemberExists: func(s1, s2 string) (bool, error) {
				return false, nil
			},
			mockIsRoleIDExists: func(u uint) (bool, error) {
				return true, nil
			},
			mockAddMember: func(cm entities.CompanyMembers) error {
				return nil
			},
			memberRequest: entities.CompanyMembers{
				CompanyID: "gdywqg6qgwu8647ddhsvy82",
				MemberID:  "djchs638i686hld9ph83hd2",
				RoleID:    5,
				StatusID:  3,
				Salary:    9000000,
			},
			isMemberExist: false,
			isRoleExist:   true,
		},
		{
			name: "Fail",
			mockIsMemberExists: func(s1, s2 string) (bool, error) {
				return true, nil
			},
			mockIsRoleIDExists: func(u uint) (bool, error) {
				return true, nil
			},
			mockAddMember: func(cm entities.CompanyMembers) error {
				return nil
			},
			memberRequest: entities.CompanyMembers{
				CompanyID: "gdywqg6qgwu8647ddhsvy82",
				MemberID:  "djchs638i686hld9ph83hd2",
				RoleID:    5,
				StatusID:  3,
				Salary:    9000000,
			},
			isMemberExist: true,
			isRoleExist:   true,
		},
		{
			name: "Fail",
			mockIsMemberExists: func(s1, s2 string) (bool, error) {
				return false, nil
			},
			mockIsRoleIDExists: func(u uint) (bool, error) {
				return false, nil
			},
			mockAddMember: func(cm entities.CompanyMembers) error {
				return nil
			},
			memberRequest: entities.CompanyMembers{
				CompanyID: "gdywqg6qgwu8647ddhsvy82",
				MemberID:  "djchs638i686hld9ph83hd2",
				RoleID:    5,
				StatusID:  3,
				Salary:    9000000,
			},
			isMemberExist: false,
			isRoleExist:   false,
		},
	}

	for _, test := range tests {

		if !test.isMemberExist && test.isRoleExist {
			adapter.EXPECT().IsMemberExists(gomock.Any(), gomock.Any()).DoAndReturn(test.mockIsMemberExists).AnyTimes().Times(1)
			adapter.EXPECT().IsRoleIDExists(gomock.Any()).DoAndReturn(test.mockIsRoleIDExists).AnyTimes().Times(1)
			adapter.EXPECT().AddMember(gomock.Any()).DoAndReturn(test.mockAddMember).AnyTimes().Times(1)
		} else if test.isMemberExist {
			adapter.EXPECT().IsMemberExists(gomock.Any(), gomock.Any()).DoAndReturn(test.mockIsMemberExists).AnyTimes().Times(1)
		} else if !test.isRoleExist {
			adapter.EXPECT().IsMemberExists(gomock.Any(), gomock.Any()).DoAndReturn(test.mockIsMemberExists).AnyTimes().Times(1)
			adapter.EXPECT().IsRoleIDExists(gomock.Any()).DoAndReturn(test.mockIsRoleIDExists).AnyTimes().Times(1)
		}
		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.AddMember(test.memberRequest)
		if !(!test.isMemberExist && test.isRoleExist) {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestLogintoCompany(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name                     string
		mockGetCompanyIDFromName func(string) (string, error)
		mockLogintoCompany       func(string, string) (entities.LogintoCompanyUsecase, error)
		companyUsernameRequest   string
		memberIDRequest          string
		wantError                bool
		expectedResult           entities.LogintoCompanyUsecase
	}{
		{
			name: "Success",
			mockGetCompanyIDFromName: func(s string) (string, error) {
				return "gkdacash7839567dghdjksh89", nil
			},
			mockLogintoCompany: func(s1, s2 string) (entities.LogintoCompanyUsecase, error) {
				return entities.LogintoCompanyUsecase{
					CompanyID:   s1,
					RoleID:      4,
					Permisssion: "SEMI-ROOT",
				}, nil
			},
			companyUsernameRequest: "itscbfdsjlhkjbdfj",
			memberIDRequest:        "cgdkgjas678vxhuso836986472bdue",
			wantError:              false,
			expectedResult: entities.LogintoCompanyUsecase{
				CompanyID:   "gkdacash7839567dghdjksh89",
				RoleID:      4,
				Permisssion: "SEMI-ROOT",
			},
		},
		{
			name: "Fail",
			mockGetCompanyIDFromName: func(s string) (string, error) {
				return "", nil
			},
			mockLogintoCompany: func(s1, s2 string) (entities.LogintoCompanyUsecase, error) {
				return entities.LogintoCompanyUsecase{
					CompanyID:   s1,
					RoleID:      4,
					Permisssion: "SEMI-ROOT",
				}, nil
			},
			companyUsernameRequest: "itscbfdsjlhkjbdfj",
			memberIDRequest:        "cgdkgjas678vxhuso836986472bdue",
			wantError:              true,
			expectedResult: entities.LogintoCompanyUsecase{
				CompanyID:   "gkdacash7839567dghdjksh89",
				RoleID:      4,
				Permisssion: "SEMI-ROOT",
			},
		},
		{
			name: "Success",
			mockGetCompanyIDFromName: func(s string) (string, error) {
				return "dssdsdfasadfsfsafsa", nil
			},
			mockLogintoCompany: func(s1, s2 string) (entities.LogintoCompanyUsecase, error) {
				return entities.LogintoCompanyUsecase{
					CompanyID:   s1,
					RoleID:      4,
					Permisssion: "SEMI-ROOT",
				}, nil
			},
			companyUsernameRequest: "itscbfdsjlhkjbdfj",
			memberIDRequest:        "cgdkgjas678vxhuso836986472bdue",
			wantError:              false,
			expectedResult: entities.LogintoCompanyUsecase{
				CompanyID:   "dssdsdfasadfsfsafsa",
				RoleID:      4,
				Permisssion: "SEMI-ROOT",
			},
		},
	}

	for _, test := range tests {

		if test.wantError {
			adapter.EXPECT().GetCompanyIDFromName(gomock.Any()).DoAndReturn(test.mockGetCompanyIDFromName).AnyTimes().Times(1)
		} else {
			adapter.EXPECT().GetCompanyIDFromName(gomock.Any()).DoAndReturn(test.mockGetCompanyIDFromName).AnyTimes().Times(1)
			adapter.EXPECT().LogintoCompany(gomock.Any(), gomock.Any()).DoAndReturn(test.mockLogintoCompany).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		res, err := regUsecase.LogintoCompany(test.companyUsernameRequest, test.memberIDRequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, test.expectedResult, res)
		}
	}
}

func TestInsertIntoClients(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name                     string
		mockGetInsertIntoClients func(entities.Clients) error
		clientRequest            entities.Clients
		wantError                bool
	}{
		{
			name: "Success",
			mockGetInsertIntoClients: func(c entities.Clients) error {
				return nil
			},
			clientRequest: entities.Clients{
				ClientID:  "jhcgskdah643786vdi36487",
				CompanyID: "cgjd2476dbo8b63478g367vd",
			},
			wantError: false,
		},
		{
			name: "Fail",
			mockGetInsertIntoClients: func(c entities.Clients) error {
				return nil
			},
			clientRequest: entities.Clients{
				ClientID:  "",
				CompanyID: "cgjd2476dbo8b63478g367vd",
			},
			wantError: true,
		},
		{
			name: "Fail",
			mockGetInsertIntoClients: func(c entities.Clients) error {
				return nil
			},
			clientRequest: entities.Clients{
				ClientID:  "jhcgskdah643786vdi36487",
				CompanyID: "",
			},
			wantError: true,
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().InsertIntoClients(gomock.Any()).DoAndReturn(test.mockGetInsertIntoClients).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.InsertIntoClients(test.clientRequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestAssignProblemToEmployee(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name                        string
		mockAssignProblemToEmployee func(string, uint) error
		employeeIDRequest           string
		problemID                   uint
		wantError                   bool
	}{
		{
			name: "Success",
			mockAssignProblemToEmployee: func(s string, u uint) error {
				return nil
			},
			employeeIDRequest: "kajsgfg783o7hfow743wewd343",
			problemID:         23,
			wantError:         false,
		},
		{
			name: "Fail",
			mockAssignProblemToEmployee: func(s string, u uint) error {
				return nil
			},
			employeeIDRequest: "",
			problemID:         23,
			wantError:         true,
		},
		{
			name: "Success",
			mockAssignProblemToEmployee: func(s string, u uint) error {
				return nil
			},
			employeeIDRequest: "kajsgfg783o7hwigew42fow743wewd343",
			problemID:         21,
			wantError:         false,
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().AssignProblemToEmployee(gomock.Any(), gomock.Any()).DoAndReturn(test.mockAssignProblemToEmployee).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.AssignProblemToEmployee(test.employeeIDRequest, test.problemID)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestApplyforLeave(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name              string
		mockApplyforLeave func(entities.Leaves) error
		leaveRequest      entities.Leaves
		wantError         bool
	}{
		{
			name: "Success",
			mockApplyforLeave: func(l entities.Leaves) error {
				return nil
			},
			leaveRequest: entities.Leaves{
				EmployeeID:  "ahdskch738dfh3478",
				CompanyID:   "hbkei7837gfg3746g83",
				Description: "avs,dhslbduer",
				Date:        time.Now(),
			},
			wantError: false,
		},
		{
			name: "Fail",
			mockApplyforLeave: func(l entities.Leaves) error {
				return nil
			},
			leaveRequest: entities.Leaves{
				EmployeeID:  "",
				CompanyID:   "hbkei7837gfg3746g83",
				Description: "avs,dhslbduer",
				Date:        time.Now(),
			},
			wantError: true,
		},
		{
			name: "Fail",
			mockApplyforLeave: func(l entities.Leaves) error {
				return nil
			},
			leaveRequest: entities.Leaves{
				EmployeeID:  "ahdskch738dfh3478",
				CompanyID:   "",
				Description: "avs,dhslbduer",
				Date:        time.Now(),
			},
			wantError: true,
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().ApplyforLeave(gomock.Any()).DoAndReturn(test.mockApplyforLeave).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.ApplyforLeave(test.leaveRequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestPostJob(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)

	tests := []struct {
		name           string
		mockPostJob    func(entities.Address, entities.Jobs) error
		addressRequest entities.Address
		jobsRequest    entities.Jobs
		wantError      bool
	}{
		{
			name: "Success",
			mockPostJob: func(a entities.Address, j entities.Jobs) error {
				return nil
			},
			addressRequest: entities.Address{
				StreetNo:   23,
				StreetName: "asbjhcd",
				PinNo:      455643,
				District:   "calicut",
				State:      "kerala",
				Nation:     "India",
			},
			jobsRequest: entities.Jobs{
				CompanyID:      "cgekwv6437v76345fd",
				Role:           "cdwvkvk",
				Vacancy:        4,
				Description:    "vcwhgvjefygwvkegfu",
				MinExperiance:  1,
				MinExpectedCTC: 7.9,
				MaxExpectedCTC: 10.9,
				IsRemote:       false,
			},
			wantError: false,
		},
		{
			name: "Fail",
			mockPostJob: func(a entities.Address, j entities.Jobs) error {
				return nil
			},
			addressRequest: entities.Address{
				StreetNo:   23,
				StreetName: "asbjhcd",
				PinNo:      455643,
				District:   "calicut",
				State:      "kerala",
				Nation:     "India",
			},
			jobsRequest: entities.Jobs{
				CompanyID:      "",
				Role:           "cdwvkvk",
				Vacancy:        4,
				Description:    "vcwhgvjefygwvkegfu",
				MinExperiance:  1,
				MinExpectedCTC: 7.9,
				MaxExpectedCTC: 10.9,
				IsRemote:       false,
			},
			wantError: true,
		},
		{
			name: "Success",
			mockPostJob: func(a entities.Address, j entities.Jobs) error {
				return nil
			},
			addressRequest: entities.Address{
				StreetNo:   23,
				StreetName: "asbjhcd",
				PinNo:      455643,
				District:   "calicut",
				State:      "kerala",
				Nation:     "India",
			},
			jobsRequest: entities.Jobs{
				CompanyID:      "cgekwv6437v76345fd",
				Role:           "cdwvkvk",
				Vacancy:        4,
				Description:    "vcwhgvjefygwvkegfu",
				MinExperiance:  1,
				MinExpectedCTC: 7.9,
				MaxExpectedCTC: 10.9,
				IsRemote:       false,
			},
			wantError: false,
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().PostJob(gomock.Any(), gomock.Any()).DoAndReturn(test.mockPostJob).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.PostJob(test.addressRequest, test.jobsRequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestApplyJob(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_adapters.NewMockCompanyAdapterInterfaces(ctrl)
	tests := []struct {
		name                    string
		mockInsertResumetoMinio func(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error
		mockApplyforJob         func(entities.JobApplications) error
		applicationRequest      entities.JobApplications
		wantError               bool
	}{
		{
			name: "Success",
			mockInsertResumetoMinio: func(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error {
				return nil
			},
			mockApplyforJob: func(ja entities.JobApplications) error {
				return nil
			},
			applicationRequest: entities.JobApplications{
				UserID: "hcbwlhej7284ofo8454",
				JobID:  "chfg4675vc28279ovro8724",
				Name:   "Akshay",
				Email:  "bdjkhq@gmail.com",
				Mobile: "6790987678",
				AddressofApplicant: entities.Address{
					StreetNo:   23,
					StreetName: "asbjhcd",
					PinNo:      455643,
					District:   "calicut",
					State:      "kerala",
					Nation:     "India",
				},
				HighestEducation: "cgwvlcryuew",
				Nationality:      "cewkllfbvw",
				Experiance:       2,
				CurrentCTC:       23.5,
			},
			wantError: false,
		},
		{
			name: "Fail",
			mockInsertResumetoMinio: func(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error {
				return nil
			},
			mockApplyforJob: func(ja entities.JobApplications) error {
				return nil
			},
			applicationRequest: entities.JobApplications{
				UserID: "",
				JobID:  "chfg4675vc28279ovro8724",
				Name:   "Akshay",
				Email:  "bdjkhq@gmail.com",
				Mobile: "6790987678",
				AddressofApplicant: entities.Address{
					StreetNo:   23,
					StreetName: "asbjhcd",
					PinNo:      455643,
					District:   "calicut",
					State:      "kerala",
					Nation:     "India",
				},
				HighestEducation: "cgwvlcryuew",
				Nationality:      "cewkllfbvw",
				Experiance:       2,
				CurrentCTC:       23.5,
			},
			wantError: true,
		},
		{
			name: "Fail",
			mockInsertResumetoMinio: func(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error {
				return nil
			},
			mockApplyforJob: func(ja entities.JobApplications) error {
				return nil
			},
			applicationRequest: entities.JobApplications{
				UserID: "hcbwlhej7284ofo8454",
				JobID:  "",
				Name:   "Akshay",
				Email:  "bdjkhq@gmail.com",
				Mobile: "6790987678",
				AddressofApplicant: entities.Address{
					StreetNo:   23,
					StreetName: "asbjhcd",
					PinNo:      455643,
					District:   "calicut",
					State:      "kerala",
					Nation:     "India",
				},
				HighestEducation: "cgwvlcryuew",
				Nationality:      "cewkllfbvw",
				Experiance:       2,
				CurrentCTC:       23.5,
			},
			wantError: true,
		},
	}

	for _, test := range tests {

		if !test.wantError {
			adapter.EXPECT().InsertResumetoMinio(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(test.mockInsertResumetoMinio).AnyTimes().Times(1)
			adapter.EXPECT().ApplyforJob(gomock.Any()).DoAndReturn(test.mockApplyforJob).AnyTimes().Times(1)
		}

		regUsecase := usecases.NewCompanyUseCases(adapter)

		err := regUsecase.ApplyJob(test.applicationRequest)
		if test.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
