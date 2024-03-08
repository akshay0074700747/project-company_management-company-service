package db

import (
	"fmt"
	"log"

	"github.com/akshay0074700747/project-company_management-company-service/config"
	"github.com/akshay0074700747/project-company_management-company-service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectDB(cfg config.Config) *gorm.DB {

	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBhost, cfg.DBuser, cfg.DBname, cfg.DBport, cfg.DBpassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("cannot connect to the db ", err)
	}

	db.AutoMigrate(&entities.CompanyTypes{})
	db.AutoMigrate(&entities.Permissions{})
	db.AutoMigrate(&entities.Credentials{})
	db.AutoMigrate(&entities.CompanyAddress{})
	db.AutoMigrate(&entities.CompanyRoles{})
	db.AutoMigrate(&entities.CompanyMembers{})
	db.AutoMigrate(&entities.CompanyEmail{})
	db.AutoMigrate(&entities.CompanyPhone{})
	db.AutoMigrate(&entities.Owners{})
	db.AutoMigrate(&entities.MemberStatus{})
	db.AutoMigrate(&entities.Problems{})
	db.AutoMigrate(&entities.Visitors{})
	db.AutoMigrate(&entities.Clients{})
	db.AutoMigrate(&entities.ClientsWithProjects{})
	db.AutoMigrate(&entities.CompanyPolicies{})
	db.AutoMigrate(&entities.PayRoll{})
	db.AutoMigrate(&entities.Leaves{})
	db.AutoMigrate(&entities.Address{})
	db.AutoMigrate(&entities.Jobs{})
	db.AutoMigrate(&entities.JobApplications{})
	db.AutoMigrate(&entities.ScheduledInterviews{})

	return db
}

func ConnectMinio(cfg config.Config) *minio.Client {
	minioClient, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
