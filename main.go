package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bejaneps/csv-webapp/crud"
	"github.com/bejaneps/csv-webapp/models"

	"github.com/bejaneps/csv-webapp/auth"

	"github.com/bejaneps/csv-webapp/handlers"
	"github.com/gin-gonic/gin"
)

var (
	logger   *os.File
	recovery *os.File
	err      error
)

var iniFile = flag.String("ini", "conf/conf.ini", "a path to ini file.")

func init() {
	flag.Parse()

	err = crud.ParseINI(*iniFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO]: Downloading ftp files\n")
	err := crud.DownloadFTPFiles()
	if err != nil {
		log.Fatal(err)
	}

	logger, err = os.Create("logs/log.txt")
	if err != nil {
		log.Fatal(err)
	}

	recovery, err = os.Create("logs/recovery.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = crud.ParseTemplates()
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode("debug")
}

func main() {
	router := gin.New()
	router.Use(gin.LoggerWithWriter(logger), gin.RecoveryWithWriter(recovery))

	router.GET("/", handlers.IndexHandler)
	router.GET("/login", handlers.LoginHandler)
	router.GET("/dashboard", handlers.DashboardHandler)
	router.GET("/report", handlers.ReportHandler)
	router.GET("/logout", handlers.LogoutHandler)

	var server = &http.Server{
		Addr:    models.Port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

	if err := auth.CloseMongoClient(); err != nil {
		log.Fatal(err)
	}

	logger.Close()
	recovery.Close()
}
