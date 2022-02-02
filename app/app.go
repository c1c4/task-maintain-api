package app

import (
	"api/app/config"
	"api/app/database"
	"api/app/database/migration"
	"api/app/middleware"
	"api/app/repositories"
	"api/app/routers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	ServerReady chan bool
}

func init() {
	var err error
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		err = godotenv.Load(os.ExpandEnv("./../../.env"))
		os.Setenv("ENV", "TEST")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Println(err)
	}
}

func (app *App) StartApp() {
	config.LoadEnv()

	switch config.ENV {
	case "PROD":
		gin.SetMode(gin.ReleaseMode)
	case "TEST":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	database.Connect()
	migration.AutoMigration()

	repositories.UserRepo.Init()
	repositories.TaskRepo.Init()

	router := gin.New()
	router.Use(middleware.Logger())

	routers.InitializeRoutes(router)

	router.SetTrustedProxies(nil)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// create an err channel to watch for server errors
	errChan := make(chan error)
	go func(errChan chan error) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(errChan)

	if app.ServerReady != nil {
		app.ServerReady <- true
	}

	// create a quit channel to watch for SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// non-blocking select to watch for errors or quit signal
	select {

	// real error logs and exits immediately
	case err := <-errChan:
		log.Fatalln(err)

	// sigint shuts the server down gracefully
	case <-quit:

		// log that we are starting the shut down
		log.Println("Shutting down server after requests finish...")

		// add a backup timeout to the context so if requests
		// don't finish in time they are cut short
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		// tell the server to shut down with the backup timeout context
		if err := server.Shutdown(ctx); err != nil {

			// if there's an error just shut down immediately
			log.Fatalln(err)
		}

		// log that the server is shutting down now
		// hopefully with all the requests finished
		log.Println("Server exiting now")
	}
}
