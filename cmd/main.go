package main

import (
	"log"

	"github.com/akshay0074700747/project-company_management-company-service/config"
	injectdependency "github.com/akshay0074700747/project-company_management-company-service/injectDependency"
)

func main() {

	config, err := config.LoadConfigurations()
	if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	injectdependency.Initialize(config).Start(":50003")

}
