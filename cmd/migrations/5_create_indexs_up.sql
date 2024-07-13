CREATE INDEX indexBrand ON car_catalog (brand);
CREATE INDEX indexModel ON car_catalog (model);
CREATE INDEX indexYear ON car_catalog (year_issue);
CREATE INDEX indexHolder ON car_catalog (holder);
CREATE INDEX indexBrandName ON brand (lower(name));
CREATE INDEX indexModelName ON model (lower(name), brand);
CREATE INDEX indexFullName ON holder (lower(name), lower(surname), lower(patronymic));