package cron

import (
	"time"

	"github.com/akshay0074700747/project-company_management-company-service/helpers"
	"gorm.io/gorm"
)

type Cron struct {
	DB *gorm.DB
}

func NewCron(db *gorm.DB) *Cron {
	return &Cron{
		DB: db,
	}
}

func (cron *Cron) Run() {

	ticker := time.NewTicker(time.Hour * 24)

	for range ticker.C {
		query := "UPDATE credentials SET is_payed = false WHERE next_payment_at < $1"
		if err := cron.DB.Exec(query, time.Now()).Error; err != nil {
			helpers.PrintErr(err, "error happened at cron")
		}
	}
}
