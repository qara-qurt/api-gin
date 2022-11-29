package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/qara-qurt/api-gin/pkg/handler"
	"github.com/qara-qurt/api-gin/pkg/repository"
	"github.com/qara-qurt/api-gin/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	gin "github.com/qara-qurt/api-gin"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitCinfig(); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatal(err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	//HTTP Server
	srv := new(gin.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatal(err.Error())
	}
}

func InitCinfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
