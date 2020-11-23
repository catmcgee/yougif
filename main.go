package main

import (
	"catmcgee/config"
	"catmcgee/controller"
	"catmcgee/repository"
	"catmcgee/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	if err := config.Load(); err != nil {
		logrus.Fatal(err)
	}

	host := config.Configuration.Db.Host
	port := config.Configuration.Db.Port
	user := config.Configuration.Db.User
	password := config.Configuration.Db.Password
	dbname := config.Configuration.Db.DatabaseName

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}

	repository.SetDatabase(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	migrations, err := migrate.NewWithDatabaseInstance("file://migrations", "posgtres", driver)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Println(migrations.Up())

	go func() {
		service.GetYouTubeVideos(config.Configuration.YouTube.Region, config.Configuration.YouTube.TotalNumbersOfVideosToDownload)
	}()

	go func() {
		service.StartProcessingVideos(config.Configuration.Processing.NumberOfWorkers)
	}()

	router := gin.Default()
	router.GET("/search", controller.Search)
	router.GET("/images/:id", controller.GetImage)
	logrus.Fatal(router.Run(fmt.Sprintf(":%d", config.Configuration.Server.Port)))
}
