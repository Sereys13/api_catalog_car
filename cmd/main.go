package main

import (
	"api_catalog_car/internal/api"
	"api_catalog_car/internal/config"
	"api_catalog_car/internal/database"
	"api_catalog_car/internal/migration"
	"embed"
	"log"

	"api_catalog_car/pkg/logging"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	logging.CreateLogger(&cfg.Logs)
	logger := logging.GetLogger()
	logger.Info("Load config")
	logger.Info("Creat logger")
	ctx := context.Background()

	logger.Info("Migrations database Postgres")
	dbm, err := database.InitDbMigration(&cfg.Storage)
	if err != nil {
		logger.Fatal("Failed connection database error = ", err)
	}

	migrator, err := migration.MustGetNewMigrator(MigrationsFS, "migrations")
	if err != nil {
		logger.Fatal("Failed migration database error = ", err)
	}
	err = migrator.ApplyMigrations(dbm)
	if err != nil {
		logger.Fatal("Failed migration database error = ", err)
	}
	defer dbm.Close()

	logger.Info("Connection database Postgres")
	db, err := database.InitDbConnect(ctx, &cfg.Storage)
	if err != nil {
		logger.Fatal("Failed connection database error = ", err)
	}
	defer db.Close()

	a := api.NewApi(ctx, db, logger, cfg.UrlApiCarInfo)

	logger.Info("Create router")
	r := mux.NewRouter()
	a.Routes(r)

	logger.Info("Listen tcp")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	if err != nil {
		logger.Fatal("Failed create listen error = ", err)
	}

	server := http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("Запуск веб-сервера на http: %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
