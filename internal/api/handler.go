package api

import (
	"api_catalog_car/internal/database"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *Api) PageCatalogGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	var catalog *database.PageCatalog
	var err error
	p := 0
	if _, ok := params["p"]; ok {
		p, err = strconv.Atoi(params["p"][0])
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	vars := mux.Vars(r)
	idObj := database.IdObject{IdBrand: vars["brand"], IdModel: vars["model"]}
	_, okFilter := params["year"]
	_, okSort := params["sort"]
	if okFilter && okSort {
		catalog, err = a.db.IssuanceCatalogSort(a.ctx,&idObj,params["year"][0],params["sort"][0],p)
	} else if okFilter {
		catalog, err = a.db.IssuanceCatalogSort(a.ctx,&idObj,params["year"][0],"",p)
	} else if okSort {
		catalog, err = a.db.IssuanceCatalogSort(a.ctx,&idObj,"",params["sort"][0],p)
	} else {
		catalog, err = a.db.IssuanceCatalog(a.ctx, &idObj, p)
	}
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	catalogJson, err := json.Marshal(catalog)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(catalogJson)

}

func (a *Api) PageCatalogDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := a.db.DeleteItemsCatalog(a.ctx, vars["id"])
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PutCatalogRegNum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var urn database.UpdateRegNum
	err = json.Unmarshal(body, &urn)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	err = a.db.UpdateItemsRegNum(a.ctx, vars["id"], &urn)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PutCatalogBrandModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ubm database.UpdateBrandModel
	err = json.Unmarshal(body, &ubm)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	err = a.db.UpdateItemsBrand(a.ctx, vars["id"], &ubm)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PutCatalogYear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var uy database.UpdateYear
	err = json.Unmarshal(body, &uy)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	err = a.db.UpdateItemsYear(a.ctx, vars["id"], &uy)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PutCatalogHolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var uh database.UpdateHolder
	err = json.Unmarshal(body, &uh)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	err = a.db.UpdateItemsHolder(a.ctx, vars["id"], &uh)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) AddItemsCatalog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type arrRegNum struct {
		RegNums []string `json:"regNums"`
	}
	var regNums arrRegNum
	err = json.Unmarshal(body, &regNums)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var resp *http.Response
	reg, err := regexp.Compile("[A-Z][0-9]{3}[A-Z]{2}[0-9]{2}[0-9]?")
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, el := range regNums.RegNums { 
		if !reg.MatchString(el){
			a.logger.Info(errors.New("Number "+el+ " no validation"))
		    w.WriteHeader(http.StatusBadRequest)
			continue
		}
		resp, err = http.Get(a.urlApiCarInfo + el)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var Item database.ItemsCatalog
		err = json.Unmarshal(body, &Item)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = a.db.AddItemsHolder(a.ctx, &Item)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) PageBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brands, err := a.db.IssuanceBrand(a.ctx)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	brandsJson, err := json.Marshal(brands)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(brandsJson)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PageModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	models, err := a.db.IssuanceModel(a.ctx, vars["brand"])
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	modelsJson, err := json.Marshal(models)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(modelsJson)
	w.WriteHeader(http.StatusOK)
}
