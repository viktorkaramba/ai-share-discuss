package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	Backend "PlaylistsSynchronizer.Backend"
	"PlaylistsSynchronizer.Backend/configs"
	"PlaylistsSynchronizer.Backend/pkg/handlers"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"PlaylistsSynchronizer.Backend/pkg/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// @title Playlist Synchronizer API
// @version 1.0
// @description API Server for Playlist Synchronizer Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load("Backend/.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	configs.LoadAuthConfig()

	db, err := repositories.NewPostgresDB(repositories.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repositories.NewRepository(db)
	service := services.NewService(repos)
	handler := handlers.NewHandler(service)

	srv := new(Backend.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("Backend/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
