CREATE TABLE IF NOT EXISTS car_catalog (
	id SERIAL PRIMARY KEY,
	regNum VARCHAR(10) NOT NULL,
	brand INT REFERENCES brand(id),
	model INT REFERENCES model(id),
	year_issue VARCHAR(4),
	holder INT REFERENCES holder(id),
	delete_status BOOLEAN DEFAULT FALSE
);