package application

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rajesh_bond/production/config"
	"github.com/rajesh_bond/production/database"
)

type App struct {
	DB     *sql.DB
	Config *config.Config
}

func NewApp() *App {
	cfg := config.Load()
	db := database.NewDB(cfg)
	return &App{
		DB:     db,
		Config: cfg,
	}
}

func (a *App) Start() error {
	// defer a.DB.Close()
	r := NewRouter(a)

	fmt.Println("Server running on:", a.Config.APPPORT)

	return http.ListenAndServe(":"+a.Config.APPPORT, r)

}
