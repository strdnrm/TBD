BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE usertg (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (), 
	username VARCHAR(255) UNIQUE
);

CREATE TABLE product (
	id UUID PRIMARY KEY  DEFAULT uuid_generate_v4 (), 
	name VARCHAR(255)
);

CREATE TABLE buy_list (
	id UUID DEFAULT uuid_generate_v4 (), 
	user_id UUID REFERENCES usertg(id),
	product_id UUID REFERENCES product(id),
	weight REAL,
	but_time TIMESTAMP
);

CREATE TABLE fridge (
	user_id UUID REFERENCES usertg(id),
	product_id UUID REFERENCES product(id),
	opened BOOLEAN,
	expire_date DATE,
	status VARCHAR(255),
	use_date DATE
);

END;