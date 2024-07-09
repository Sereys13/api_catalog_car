package api

import (
	"api_catalog_car/internal/database"
	"api_catalog_car/pkg/logging"
	"context"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Api struct {
	ctx    context.Context
	db     *database.DB
	logger *logging.Logger
	urlApiCarInfo string
}

func NewApi(ctx context.Context, db *pgxpool.Pool, logger *logging.Logger, otherApi string) *Api {
	return &Api{ctx: ctx, db: database.NewDataBase(db, logger), logger: logger, urlApiCarInfo: otherApi}
}

func (a *Api) Routes(r *mux.Router) {
	//Хэндлер выдачи общего каталога
	r.HandleFunc("/catalog", a.PageCatalogGet).Methods("GET")
	//хэндлер добавления записей в каталоге (в данном хэндлере необходимо добавить url  к api)
	r.HandleFunc("/catalog", a.AddItemsCatalog).Methods("POST")
	//Хэндлер удаление записи в каталоге
	r.HandleFunc("/catalog/{id:[0-9]+}", a.PageCatalogDelete).Methods("DELETE")
	//Хэндлеры изменения поля/полей в записи в каталоге
	r.HandleFunc("/catalog/{id:[0-9]+}", a.PutCatalog).Methods("PUT")
	//Хэндлеры изменения поля/полей в записи в каталоге
	r.HandleFunc("/filters", a.GetFilters).Methods("GET")
}
