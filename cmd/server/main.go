package main

import (
	"log"

	application "github.com/rajesh_bond/production/Application"
	"github.com/rajesh_bond/production/internal/common/utils"

	_ "github.com/rajesh_bond/production/docs"
)

//	@title			Production API
//	@version		1.0
//	@description	Mold Management SaaS API
//	@host			localhost:8080
//	@BasePath		/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	utils.InitValidator()
	app := application.NewApp()
	defer app.DB.Close()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
