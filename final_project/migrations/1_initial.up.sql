BEGIN

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
	id UUID PRIMARY KEY ,
	login VARCHAR(128) UNIQUE,
	email VARCHAR(128) UNIQUE,
    password VARCHAR(128),
    name VARCHAR(128),
    surname VARCHAR(128),
    is_admin BOOLEAN
); 

CREATE TABLE plane(
    id UUID PRIMARY KEY ,
    number_of_seats INT,
    model VARCHAR(128)
);

CREATE TABLE route(
    id UUID PRIMARY KEY ,
    source VARCHAR(128),
    destination VARCHAR(128)
);

CREATE TABLE flight(
    id UUID PRIMARY KEY ,
    plane_id UUID NOT NULL REFERENCES plane(id),
    route_id UUID NOT NULL REFERENCES route(id),
    departure_time TIMESTAMP,
    arrival_time TIMESTAMP,
    available_seats INT,
    transfer_flight_id UUID REFERENCES flight(id)
);

CREATE TABLE seat(
    id UUID PRIMARY KEY,
    class VARCHAR(30),
    seat_number VARCHAR(5),
    user_id UUID REFERENCES users(id),
    flight_id UUID NOT NULL REFERENCES flight(id)
);

CREATE TABLE ticket(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    flight_id UUID NOT NULL REFERENCES route(id),
    price DECIMAL,
    seat_id UUID NOT NULL REFERENCES seat(id)    
);

CREATE TABLE user_tickets(
	user_id UUID NOT NULL REFERENCES users(id),
    ticket_id UUID NOT NULL REFERENCES ticket(id)
);

END;