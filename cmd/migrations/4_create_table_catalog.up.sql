CREATE TABLE IF NOT EXISTS car_catalog (
	id SERIAL PRIMARY KEY,
	regNum VARCHAR(10) NOT NULL,
	brand INT REFERENCES brand(id) NOT NULL,
	model INT REFERENCES model(id) NOT NULL,
	year_issue VARCHAR(4) DEFAULT 'N/A',
	holder INT REFERENCES holder(id) NOT NULL,
	delete_status BOOLEAN DEFAULT FALSE
);