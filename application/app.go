package application

import (
	"log"
	"net/http"

	"github.com/rajesh_bond/production/config"
	"github.com/rajesh_bond/production/database"
)

type App struct {
	DB     *database.DB // ✅ updated
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

	// 🔥 Graceful shutdown handling (recommended)
	defer func() {
		if a.DB != nil {
			a.DB.Close()
		}
	}()

	r := NewRouter(a)

	log.Println("🚀 Server running on:", a.Config.APPPORT)

	return http.ListenAndServe(":"+a.Config.APPPORT, r)
}

// OLD code in case of error

// type App struct {
// 	DB     *sql.DB
// 	Config *config.Config
// }

// func NewApp() *App {
// 	cfg := config.Load()
// 	db := database.NewDB(cfg)
// 	return &App{
// 		DB:     db,
// 		Config: cfg,
// 	}
// }

// func (a *App) Start() error {
// 	// defer a.DB.Close()
// 	r := NewRouter(a)

// 	fmt.Println("Server running on:", a.Config.APPPORT)

// 	return http.ListenAndServe(":"+a.Config.APPPORT, r)

// }
