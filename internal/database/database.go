package database

import (
	"api_catalog_car/internal/config"
	"api_catalog_car/pkg/logging"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db     *pgx.Conn
	logger *logging.Logger
}

func NewDataBase(db *pgx.Conn, logger *logging.Logger) *DB {
	return &DB{db: db, logger: logger}
}

func InitDbConnect(ctx context.Context, cfg *config.Storage) (db *pgx.Conn, err error) {
	var connectDatabase string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.Sslmode)
	db, err = pgx.Connect(ctx, connectDatabase)
	return
}

func InitDbMigration(cfg *config.Storage) (db *sql.DB, err error) {
	var connectDatabase string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.Sslmode)
	db, err = sql.Open("pgx", connectDatabase)
	return
}
