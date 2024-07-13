package database

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
)

func (d *DB) IssuanceCatalog(ctx context.Context, page, limit int) (pc PageCatalog, err error) {
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		d.logger.Error(err)
		return 
	}

	defer conn.Release()
	row, err := conn.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
			CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
			JOIN brand b ON cc.brand = b.id 
			JOIN model m ON cc.model = m.id 
			JOIN holder h ON cc.holder = h.id 
			WHERE cc.delete_status = false AND cc.id > $1 ORDER BY cc.id ASC LIMIT $2`, page, limit)
	if err != nil {
		d.logger.Info(err)
		return
	}

	defer row.Close()
	pc.Catalog = make([]Catalog, 0, limit)
	var c Catalog
	var yearNull sql.NullInt64
	for row.Next(){
		err = row.Scan(&c.Id, &c.RegNum, &c.Brand, &c.Model, &yearNull, &c.FullName)
		if err != nil {
			d.logger.Info(err)
			return
		}

		if yearNull.Valid{
			c.Year = int(yearNull.Int64)
		} else {
			c.Year = 0
		}
		pc.Catalog = append(pc.Catalog, c)
	}

	if len(pc.Catalog) != 0 {
		pc.LastInd = pc.Catalog[len(pc.Catalog)-1].Id
		return pc, err
	} else {
		return 
	}
}

func (d *DB) IssuanceCatalogWithFiltr(ctx context.Context, fq *FiltrsQuery) (pc PageCatalog, err error) {
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		d.logger.Error(err)
		return 
	}

	defer conn.Release()

	var queryBrand, queryModel, queryHolder, queryYears string
	if len(fq.Brands) != 0 {
		queryBrand = " AND (b.id = " + strings.Join(fq.Brands, " OR b.id = ") + " )"
	}

	if len(fq.Models) != 0 {
		queryModel = " AND (m.id = " + strings.Join(fq.Models, " OR m.id = ") + " )"
	}

	if len(fq.Holders) != 0 {
		queryHolder = " AND (cc.holder = " + strings.Join(fq.Holders, " OR cc.holder = ") + " )"
	}

	
	if len(fq.Years) != 0 {
		if fq.Years[0] == "ravno"{
			queryYears = " AND (cc.year_issue = " + strings.Join(fq.Years[1:], " OR cc.year_issue = ") + " )"
		} else if fq.Years[0] == "do"{
			queryYears = " AND (cc.year_issue <= " + fq.Years[1] + " )"
		} else {
			if len(fq.Years) == 3{
				queryYears = " AND (cc.year_issue >= " + fq.Years[1] + " AND cc.year_issue <= " + fq.Years[2] + " )"
			} else {
				queryYears = " AND (cc.year_issue >= " + fq.Years[1] + " )"
			}
		}
	}

	sort := "AND cc.id > $1 ORDER BY cc.id ASC LIMIT $2"
	row, err := conn.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
			CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
			JOIN brand b ON cc.brand = b.id 
			JOIN model m ON cc.model = m.id 
			JOIN holder h ON cc.holder = h.id 
			WHERE cc.delete_status = false` + queryBrand + queryModel + queryHolder + queryYears + sort, 
			fq.Page, fq.Limit)
	if err != nil {
		d.logger.Info(err)
		return
	}

	defer row.Close()
	pc.Catalog = make([]Catalog, 0, fq.Limit)
	var c Catalog
	var yearNull sql.NullInt64
	for row.Next(){
		err = row.Scan(&c.Id, &c.RegNum, &c.Brand, &c.Model, &yearNull, &c.FullName)
		if err != nil {
			d.logger.Info(err)
			return
		}

		if yearNull.Valid{
			c.Year = int(yearNull.Int64)
		} else {
			c.Year = 0
		}
		pc.Catalog = append(pc.Catalog, c)
	}

	if len(pc.Catalog) != 0 {
		pc.LastInd = pc.Catalog[len(pc.Catalog)-1].Id
		return pc, err
	} else {
		return 
	}
}

func (d *DB) DeleteItemsCatalog(ctx context.Context, idItems string) error {
	_, err := d.db.Exec(ctx, `UPDATE car_catalog SET delete_status = true WHERE id = $1`, idItems)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) UpdateItems(ctx context.Context, idItems string, newData *UpdateCatalog) (err error){
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		d.logger.Error(err)
		return
	}

	var qu QueryUpdate
	var nullId sql.NullInt64
	var quers []string
	if newData.Brand != ""{
		err = conn.QueryRow(ctx, `SELECT id FROM brand WHERE LOWER(name) = LOWER($1)`, newData.Brand).Scan(&nullId)
		if nullId.Valid{
			qu.Brand = int(nullId.Int64)
			nullId.Valid = false
		} else {
			err = conn.QueryRow(ctx, `INSERT INTO brand (name) VALUES ($1) RETURNING id`, capitalizeFirstLetter(newData.Brand)).Scan(&qu.Brand)
		}

		if err != nil {
			d.logger.Error(err)
			return
		}
		quers = append(quers, "brand = " + strconv.Itoa(qu.Brand))
	}

	if newData.Model != ""{
		err = conn.QueryRow(ctx, `SELECT id FROM model WHERE LOWER(name) = LOWER($1) AND brand = $2`, newData.Model, qu.Brand).Scan(&nullId)
		if nullId.Valid{
			qu.Model = int(nullId.Int64)
			nullId.Valid = false
		} else {
			err = conn.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) RETURNING id`, newData.Model, qu.Brand).Scan(&qu.Model)
		}

		if err != nil {
			d.logger.Error(err)
			return
		}
		quers = append(quers, "model = " + strconv.Itoa(qu.Model))
	}

	if newData.Name != ""{
		var nullString sql.NullString
		if newData.Patronymic != ""{
			nullString.Valid = true
			nullString.String = newData.Patronymic
		}

		if newData.Patronymic == ""{
			err = conn.QueryRow(ctx, `SELECT id FROM holder WHERE LOWER(name) = LOWER($1)
							AND LOWER(surname) = LOWER($2) AND patronymic IS NULL`, newData.Name, newData.Surname).Scan(&nullId)
		} else {
			err = conn.QueryRow(ctx, `SELECT id FROM holder WHERE LOWER(name) = LOWER($1)
							AND LOWER(surname) = LOWER($2) AND LOWER(patronymic) = LOWER($3)`, newData.Name, newData.Surname, nullString).Scan(&nullId)
		}
		if nullId.Valid{
			qu.Holder = int(nullId.Int64)
			nullId.Valid = false
		} else {
			err = conn.QueryRow(ctx, `INSERT INTO holder (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id`, newData.Name, newData.Surname, nullString).Scan(&qu.Holder)
		}

		if err != nil {
			d.logger.Error(err)
			return
		}

		quers = append(quers, "holder = " + strconv.Itoa(qu.Holder) + " ")
	}

	if newData.RegNum != ""{
		quers = append(quers, "regnum = '" + newData.RegNum + "' ")
	}

	if newData.Year != 0{
		quers = append(quers, "year_issue = " + strconv.Itoa(newData.Year) + " ")
	}

	query := "UPDATE car_catalog SET " + strings.Join(quers, ",") +" WHERE id = $1"
	_, err = conn.Exec(ctx, query, idItems)
	if err != nil {
		d.logger.Error(err)
		return
	}
	return
}

func (d *DB) AddItemsHolder(ctx context.Context, items *ItemsCatalog) error {
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		d.logger.Error(err)
		return err
	}

	defer conn.Release()
	var idBrand, idModel sql.NullInt64
	err = conn.QueryRow(ctx, `SELECT id FROM brand WHERE LOWER(name) = LOWER($1)`, items.Brand).Scan(&idBrand)
	if !idBrand.Valid{
		err = conn.QueryRow(ctx, `INSERT INTO brand (name) VALUES ($1) returning id`, capitalizeFirstLetter(items.Brand)).Scan(&idBrand)
		if err != nil{
			d.logger.Error(err)
			return err
		}

		err = conn.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, items.Model, idBrand.Int64).Scan(&idModel)
		if err != nil{
			d.logger.Error(err)
			return err
		}

	} else {
		if err != nil{
			d.logger.Error(err)
			return err
		}

		err = conn.QueryRow(ctx, `SELECT id FROM model WHERE LOWER(name) = LOWER($1) AND brand = $2`, items.Model, idBrand.Int64).Scan(&idModel)
		if !idModel.Valid{
			err = conn.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, items.Model, idBrand.Int64).Scan(&idModel)
		}

		if err != nil {
			d.logger.Error(err)
			return err
		}
	}

	var idHolder sql.NullInt64
	var sqlPatronymic sql.NullString
	var idCarCatalog sql.NullInt64
	if items.Owner.Patronymic != ""{
		err = conn.QueryRow(ctx, `SELECT id FROM holder WHERE LOWER(name) = LOWER($1) AND LOWER(surname) = LOWER($2) AND LOWER(patronymic) = LOWER($3)`, items.Owner.Name, items.Owner.Surname, items.Owner.Patronymic).Scan(&idHolder)
		sqlPatronymic.Valid = true
		sqlPatronymic.String = items.Owner.Patronymic
	} else {
		err = conn.QueryRow(ctx, `SELECT id FROM holder WHERE LOWER(name) = LOWER($1) AND LOWER(surname) = LOWER($2) AND patronymic IS NULL`, items.Owner.Name, items.Owner.Surname).Scan(&idHolder)
	} 

	if !idHolder.Valid{
		err = conn.QueryRow(ctx, `INSERT INTO holder (name, surname, patronymic) VALUES ($1, $2, $3) returning id`, capitalizeFirstLetter(items.Owner.Name), capitalizeFirstLetter(items.Owner.Surname), sqlPatronymic).Scan(&idHolder)
	}

	if err != nil {
		d.logger.Error(err)
		return err
	}

	var year sql.NullInt64
	if items.Year == 0 {
		year.Valid = false
		err = conn.QueryRow(ctx, `SELECT id FROM car_catalog WHERE regnum = $1 AND brand = $2 AND model = $3 AND holder = $4 AND year_issue IS NULL`, items.RegNum, idBrand.Int64, idModel.Int64, idHolder.Int64).Scan(&idCarCatalog)
	} else {
		year.Valid = true
		year.Int64 = int64(items.Year)
		err = conn.QueryRow(ctx, `SELECT id FROM car_catalog WHERE regnum = $1 AND brand = $2 AND model = $3 AND year_issue = $4 AND holder = $5`, items.RegNum, idBrand.Int64, idModel.Int64, items.Year, idHolder.Int64).Scan(&idCarCatalog)
	}

	if !idCarCatalog.Valid{
		_, err = conn.Exec(ctx, `INSERT INTO car_catalog (regnum, brand, model, year_issue, holder) VALUES ($1, $2, $3, $4, $5)`, items.RegNum, idBrand.Int64, idModel.Int64, year, idHolder.Int64)
	} else {
		if err != nil {
			d.logger.Error(err)
			return err
		}
		_, err = conn.Exec(ctx, `UPDATE car_catalog SET delete_status = false WHERE id = $1`, idCarCatalog.Int64)
	}

	if err != nil {
		d.logger.Error(err)
		return err
	}

	return err
}

func (d *DB) IssuanceFilters(ctx context.Context) (f Filters, err error){
	conn, err := d.db.Acquire(ctx)
	if err != nil {
		d.logger.Error(err)
		return
	}

	defer conn.Release()
	f.Brands = make(map[string]*ModelFilters, 20)
	row, err := conn.Query(ctx, `SELECT m.brand, b.name, m.id, m.name FROM model m
								JOIN brand b ON m.brand = b.id`)
	if err != nil {
		d.logger.Error(err)
		return
	}
	defer row.Close()
	
	var idBrand, idModel int
	var nameBrand, nameModel string
	var ok bool

	for row.Next(){
		err = row.Scan(&idBrand, &nameBrand, &idModel, &nameModel)
		if err != nil {
			d.logger.Error(err)
			return
		}

		if _, ok = f.Brands[nameBrand]; ok {
			f.Brands[nameBrand].Models = append(f.Brands[nameBrand].Models, Model{Id: idModel, Name: nameModel})
		} else {
			f.Brands[nameBrand] = &ModelFilters{IdBrand: idBrand, Models: make([]Model, 0, 7)}
			f.Brands[nameBrand].Models = append(f.Brands[nameBrand].Models, Model{Id: idModel, Name: nameModel})
		}
	}

	row, err = conn.Query(ctx, `SELECT id, CONCAT(surname, ' ', Name, ' ', patronymic) FROM holder`)
	if err != nil {
		d.logger.Error(err)
		return
	}

	f.Holders, err = pgx.CollectRows(row, pgx.RowToStructByPos[Holder])
	if err != nil {
		d.logger.Info(err)
		return
	}

	return
}

func capitalizeFirstLetter(s string) string {
    r := []rune(s)
    r[0] = unicode.ToUpper(r[0])
    for i := 1; i < len(r); i++ {
        r[i] = unicode.ToLower(r[i])
    }
    return string(r)
}