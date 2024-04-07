package api

import (
	"api_catalog_car/internal/database"
	"api_catalog_car/pkg/logging"
	"context"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Api struct {
	ctx    context.Context
	db     *database.DB
	logger *logging.Logger
}

func NewApi(ctx context.Context, db *pgx.Conn, logger *logging.Logger) *Api {
	return &Api{ctx: ctx, db: database.NewDataBase(db, logger), logger: logger}
}

func (a *Api) Routes(r *mux.Router) {
	//Хэндлер выдачи общего каталога
	r.HandleFunc("/catalog", a.PageCatalogGet).Methods("GET")
	//Хэндлер выдачи каталога определеного бренда автомобиля
	r.HandleFunc("/catalog/{brand:[0-9]+}", a.PageCatalogGet).Methods("GET")
	//Хэндлер выдачи каталога определеной марки бренда автомобиля
	r.HandleFunc("/catalog/{brand:[0-9]+}/{model:[0-9]+}", a.PageCatalogGet).Methods("GET")
	//хэндлер добавления записей в каталоге (в данном хэндлере необходимо добавить url  к api)
	r.HandleFunc("/catalog", a.AddItemsCatalog).Methods("POST")
	//Хэндлер удаление записи в каталоге
	r.HandleFunc("/catalog/{id:[0-9]+}", a.PageCatalogDelete).Methods("DELETE")
	//Хэндлеры изменения поля/полей в записи в каталоге
	r.HandleFunc("/catalog/regnum/{id:[0-9]+}", a.PutCatalogRegNum).Methods("PUT")
	r.HandleFunc("/catalog/brand/{id:[0-9]+}", a.PutCatalogBrandModel).Methods("PUT")
	r.HandleFunc("/catalog/model/{id:[0-9]+}", a.PutCatalogBrandModel).Methods("PUT")
	r.HandleFunc("/catalog/year/{id:[0-9]+}", a.PutCatalogYear).Methods("PUT")
	r.HandleFunc("/catalog/holder/{id:[0-9]+}", a.PutCatalogHolder).Methods("PUT")
	//хэндлер выдачи списка брендов автомобилей
	r.HandleFunc("/brand", a.PageBrand)
	//хэндлер выдачи списка марок бренда автомобилей
	r.HandleFunc("/brand/{brand:[0-9]+}", a.PageModel)

}
