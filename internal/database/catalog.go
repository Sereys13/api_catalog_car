package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"strconv"
)

func (d *DB) IssuanceCatalog(ctx context.Context, page int, filtrYear string) (*PageCatalog, error) {
	var row pgx.Rows
	var err error
	if filtrYear != "" {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
       CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
	   JOIN brand b ON cc.brand = b.id 
	   JOIN model m ON cc.model = m.id 
	   JOIN holder h ON cc.holder = h.id 
	   WHERE cc.delete_status = false AND cc.id > $1 AND cc.year_issue = $2 ORDER BY cc.id ASC LIMIT 10`, page, filtrYear)
	} else {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
       CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
	   JOIN brand b ON cc.brand = b.id 
	   JOIN model m ON cc.model = m.id 
	   JOIN holder h ON cc.holder = h.id 
	   WHERE cc.delete_status = false AND cc.id > $1 ORDER BY cc.id ASC LIMIT 10`, page)
	}
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	var pc PageCatalog
	pc.Catalog = make([]Catalog, 0, 20)
	pc.Catalog, err = pgx.CollectRows(row, pgx.RowToStructByPos[Catalog])
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	pc.LastInd = pc.Catalog[len(pc.Catalog)-1].Id
	return &pc, err
}

func (d *DB) IssuanceCatalogBrand(ctx context.Context, page int, brandId string, filtrYear string) (*PageCatalog, error) {
	var row pgx.Rows
	var err error
	if filtrYear != "" {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
       CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
	   JOIN brand b ON cc.brand = b.id 
	   JOIN model m ON cc.model = m.id 
	   JOIN holder h ON cc.holder = h.id 
	   WHERE cc.delete_status = false AND cc.id > $1 AND cc.year_issue = $2 AND cc.brand = $3 ORDER BY cc.id ASC LIMIT 10`, page, filtrYear, brandId)
	} else {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
        CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
		JOIN brand b ON cc.brand = b.id 
		JOIN model m ON cc.model = m.id 
		JOIN holder h ON cc.holder = h.id 
		WHERE cc.delete_status = false AND cc.brand = $1 AND cc.id > $2 ORDER BY cc.id ASC LIMIT 10`, brandId, page)
	}
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	var pc PageCatalog
	pc.Catalog = make([]Catalog, 0, 20)
	pc.Catalog, err = pgx.CollectRows(row, pgx.RowToStructByPos[Catalog])
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	pc.LastInd = pc.Catalog[len(pc.Catalog)-1].Id
	return &pc, err
}

func (d *DB) IssuanceCatalogModel(ctx context.Context, page int, brandId, modelId string, filtrYear string) (*PageCatalog, error) {
	var row pgx.Rows
	var err error
	if filtrYear != "" {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
          CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
	      JOIN brand b ON cc.brand = b.id 
	      JOIN model m ON cc.model = m.id 
	      JOIN holder h ON cc.holder = h.id 
	      WHERE cc.delete_status = false AND cc.id > $1 AND cc.year_issue = $2 AND cc.brand = $3 AND cc.Model = $4 ORDER BY cc.id ASC LIMIT 10`, page, filtrYear, brandId, modelId)
	} else {
		row, err = d.db.Query(ctx, `SELECT cc.id, cc.regnum, b.name, m.name, cc.year_issue, 
       CONCAT(h.surname, ' ', h.Name, ' ', h.patronymic) AS fullName FROM car_catalog cc
	   JOIN brand b ON cc.brand = b.id 
	   JOIN model m ON cc.model = m.id 
	   JOIN holder h ON cc.holder = h.id 
	   WHERE cc.delete_status = false AND cc.brand = $1 AND cc.Model = $2 AND cc.id > $3 ORDER BY cc.id ASC LIMIT 10`, brandId, modelId, page)
	}
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	var pc PageCatalog
	pc.Catalog = make([]Catalog, 0, 20)
	pc.Catalog, err = pgx.CollectRows(row, pgx.RowToStructByPos[Catalog])
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	pc.LastInd = pc.Catalog[len(pc.Catalog)-1].Id
	return &pc, err
}

func (d *DB) DeleteItemsCatalog(ctx context.Context, idItems string) error {
	_, err := d.db.Exec(ctx, `UPDATE car_catalog SET delete_status = true WHERE id = $1`, idItems)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) UpdateItemsRegNum(ctx context.Context, idItems string, urm *UpdateRegNum) error {
	_, err := d.db.Exec(ctx, `UPDATE car_catalog SET regnum = $1 WHERE id = $2`, urm.RegNum, idItems)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) UpdateItemsBrand(ctx context.Context, idItems string, ubm *UpdateBrandModel) error {
	var idModel, idBrand int
	err := d.db.QueryRow(ctx, `SELECT id FROM brand WHERE name = $1`, ubm.Brand).Scan(&idBrand)
	if err != nil && err != pgx.ErrNoRows {
		d.logger.Info(err)
		return err
	}
	if idBrand == 0 {
		err = d.db.QueryRow(ctx, `INSERT INTO brand (name) VALUES ($1) returning id`, ubm.Brand).Scan(&idBrand)
		if err != nil {
			d.logger.Info(err)
			return err
		}
		err = d.db.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, ubm.Model, idBrand).Scan(&idModel)
		if err != nil {
			d.logger.Info(err)
			return err
		}
	} else {
		err = d.db.QueryRow(ctx, `SELECT id FROM model WHERE name = $1 AND brand = $2`, ubm.Model, idBrand).Scan(&idModel)
		if err != nil && err != pgx.ErrNoRows {
			d.logger.Info(err)
			return err
		}
	}
	if idModel == 0 {
		err = d.db.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, ubm.Model, idBrand).Scan(&idModel)
		if err != nil && err != pgx.ErrNoRows {
			d.logger.Info(err)
			return err
		}
	}

	_, err = d.db.Exec(ctx, `UPDATE car_catalog SET model = $1, brand = $2 WHERE id = $3`, idModel, idBrand, idItems)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) UpdateItemsYear(ctx context.Context, idItems string, uy *UpdateYear) error {
	_, err := d.db.Exec(ctx, `UPDATE car_catalog SET year_issue = $1 WHERE id = $2`, uy.Year, idItems)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) UpdateItemsHolder(ctx context.Context, idItems string, uh *UpdateHolder) error {
	var idHolder int
	err := d.db.QueryRow(ctx, `SELECT holder FROM car_catalog WHERE id = $1`, idItems).Scan(&idHolder)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	_, err = d.db.Exec(ctx, `UPDATE holder SET name = $1, surname = $2, patronymic = $3 WHERE id = $4`, uh.Name, uh.Surname, uh.Patronymic, idHolder)
	if err != nil {
		d.logger.Info(err)
		return err
	}
	return err
}

func (d *DB) AddItemsHolder(ctx context.Context, items *ItemsCatalog) error {
	var idBrand, idModel int
	err := d.db.QueryRow(ctx, `SELECT id FROM brand WHERE name = $1`, items.Brand).Scan(&idBrand)
	if err != nil && err != pgx.ErrNoRows {
		d.logger.Info(err)
		return err
	}
	if idBrand == 0 {
		err = d.db.QueryRow(ctx, `INSERT INTO brand (name) VALUES ($1) returning id`, items.Brand).Scan(&idBrand)
		if err != nil && err != pgx.ErrNoRows {
			d.logger.Info(err)
			return err
		}
		err = d.db.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, items.Model, idBrand).Scan(&idModel)
		if err != nil {
			d.logger.Info(err)
			return err
		}
	} else {
		err = d.db.QueryRow(ctx, `SELECT id FROM model WHERE name = $1 AND brand = $2`, items.Model, idBrand).Scan(&idModel)
		if err != nil && err != pgx.ErrNoRows {
			d.logger.Info(err)
			return err
		}
	}
	if idModel == 0 {
		err = d.db.QueryRow(ctx, `INSERT INTO model (name, brand) VALUES ($1, $2) returning id`, items.Model, idBrand).Scan(&idModel)
		if err != nil {
			d.logger.Info(err)
			return err
		}
	}
	var idHolder int
	err = d.db.QueryRow(ctx, `SELECT id FROM car_catalog WHERE brand = $1 AND model = $2 AND regnum = $3`, idBrand, idModel, items.RegNum).Scan(&idHolder)
	if err == pgx.ErrNoRows {
		err = d.db.QueryRow(ctx, `INSERT INTO holder (name, surname, patronymic) VALUES ($1, $2,$3) returning id`, items.Owner.Name, items.Owner.Surname, items.Owner.Patronymic).Scan(&idHolder)
		if err != nil {
			d.logger.Info(err)
			return err
		}
		year := ""
		if items.Year == 0 {
			year = "N/A"
		} else {
			year = strconv.Itoa(items.Year)
		}
		_, err = d.db.Exec(ctx, `INSERT INTO car_catalog (regnum,brand,model,year_issue,holder) VALUES ($1,$2,$3,$4,$5)`,
			items.RegNum, idBrand, idModel, year, idHolder)
		if err != nil {
			d.logger.Info(err)
			return err
		}
	} else if err == nil {
		return errors.New("this line already exists")
	}
	return err
}

func (d *DB) IssuanceBrand(ctx context.Context) ([]Brand, error) {
	row, err := d.db.Query(ctx, `SELECT id, name FROM brand`)
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	var sliceBrand []Brand
	sliceBrand, err = pgx.CollectRows(row, pgx.RowToStructByPos[Brand])
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	return sliceBrand, err
}

func (d *DB) IssuanceModel(ctx context.Context, idBrand string) ([]Model, error) {
	row, err := d.db.Query(ctx, `SELECT id, name FROM model WHERE brand = $1`, idBrand)
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	var sliceModel []Model
	sliceModel, err = pgx.CollectRows(row, pgx.RowToStructByPos[Model])
	if err != nil {
		d.logger.Info(err)
		return nil, err
	}
	return sliceModel, err
}
