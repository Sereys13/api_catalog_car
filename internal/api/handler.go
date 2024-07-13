package api

import (
	"api_catalog_car/internal/database"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) PageCatalogGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	var err error
	p := 0
	limit := 10
	if _, ok := params["p"]; ok {
		p, err = strconv.Atoi(params["p"][0])
		if err != nil {
			a.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if _, ok := params["count"]; ok {
		limit, err = strconv.Atoi(params["count"][0])
		if err != nil {
			a.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var catalog database.PageCatalog
	if _, ok := params["filtr"]; ok {
		fq := database.FiltrsQuery{Page: p, Limit: limit}
		if brands, ok := params["brand"]; ok {
			for _, el := range brands{
				_, err = strconv.Atoi(el)
				if err != nil {
					a.logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: %s", "Некорретный параметр запроса brand")
					return
				}
			}
			fq.Brands = brands
		}

		if models, ok := params["model"]; ok {
			for _, el := range models{
				_, err = strconv.Atoi(el)
				if err != nil {
					a.logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: %s", "Некорретный параметр запроса model")
					return
				}
			}
			fq.Models = models
		}

		if holders, ok := params["holder"]; ok {
			for _, el := range holders{
				_, err = strconv.Atoi(el)
				if err != nil {
					a.logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: %s", "Некорретный параметр запроса holder")
					return
				}
			}
			fq.Holders = holders
		}

		if years, ok := params["year"]; ok {
			if years[0] == "ot" {
				if len(years) < 2 || len(years) > 3{
					err = errors.New("параметр year тип ot может содержать не больше трех значений")
				}
			} else if years[0] == "do" {
				if len(years) < 2 || len(years) > 3{
					err = errors.New("параметр year тип do может содержать не больше двух значений")
				}
			} else if years[0] != "ravno"{
				err = errors.New("неизвестный тип параметра year")
			}

			if err != nil {
				a.logger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Error: %s", "Некорретный параметр запроса year")
				return
			}

			for _, el := range years[1:]{
				_, err = strconv.Atoi(el)
				if err != nil {
					a.logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: %s", "Некорретный параметр запроса year")
					return
				}
			}
			fq.Years = years
		}

		catalog, err = a.db.IssuanceCatalogWithFiltr(a.ctx, &fq)
	} else {
		catalog, err = a.db.IssuanceCatalog(a.ctx, p, limit)
	}

	if err != nil {
		a.logger.Error(err)
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

func (a *Api) PutCatalog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data database.UpdateCatalog
	err = json.Unmarshal(body, &data)
	if err != nil {
		a.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	if (data.Model != "" && data.Brand == "") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", "Для смены модели необходимо заполнить поле brand")
		return
	} else if (data.Name == "" && data.Surname !="")  || (data.Name != "" && data.Surname =="")  || (data.Patronymic != "" && (data.Surname == "" || data.Name == "")){
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", "Для смены владельца поля имя и фамилия должны быть заполнены")
		return
	} else if data.Year != 0 && (data.Year < 1900 || data.Year > time.Now().Year()){
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", "Год выпуска не корректен")
		return
	} else if data.RegNum != ""{
		reg, err := regexp.Compile("[A-Z][0-9]{3}[A-Z]{2}[0-9]{2}[0-9]?$")
		if err != nil {
			a.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !reg.MatchString(data.RegNum){
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", "Некоректный номер для смены")
			return
		}
	}

	vars := mux.Vars(r)
	err = a.db.UpdateItems(a.ctx, vars["id"], &data)
	if err != nil {
		a.logger.Error(err)
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
	reg, err := regexp.Compile("[A-Z][0-9]{3}[A-Z]{2}[0-9]{2}[0-9]?$")
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
		resp, err = http.Get(a.urlApiCarInfo + "/info?regNum=" + el)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		if resp.StatusCode != http.StatusOK{
			a.logger.Info("Ошибка на стороннем сервере")
			w.WriteHeader(resp.StatusCode)
			return
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			a.logger.Info(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(body) == 0 {
			a.logger.Info("Автомобиля с номером " + el + " нет в базе")
			fmt.Fprintf(w, "Error: %s", "Автомобиля с номером " + el + " нет в базе")
			continue
		}

		var Item database.ItemsCatalog
		err = json.Unmarshal(body, &Item)
		if err != nil {
			a.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if Item.Year != 0 && (Item.Year < 1900 || Item.Year > time.Now().Year()){
			a.logger.Info(errors.New("No-correct year "+ el + " no validation"))
		    w.WriteHeader(http.StatusBadRequest)
			continue
		} else if Item.Brand == "" || Item.Model == "" || Item.Owner.Name == "" || Item.Owner.Surname == "" || el != Item.RegNum {
			a.logger.Info(errors.New("No-correct data "+ el + " no validation"))
		    w.WriteHeader(http.StatusBadRequest)
			continue
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

func (a *Api) GetFilters(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	filters, err := a.db.IssuanceFilters(a.ctx)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		a.logger.Info(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(filtersJSON)
}

