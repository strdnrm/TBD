BEGIN

CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	login VARCHAR(128) UNIQUE,
	email VARCHAR(128) UNIQUE,
    password VARCHAR(128),
    name VARCHAR(128),
    surname VARCHAR(128),
    is_admin BOOLEAN
); 

CREATE TABLE plane(
    id SERIAL PRIMARY KEY,
    number_of_seats INT,
    model VARCHAR(128)
);

CREATE TABLE route(
    id SERIAL PRIMARY KEY,
    source VARCHAR(128),
    destination VARCHAR(128)
);

CREATE TABLE flight(
    id SERIAL PRIMARY KEY,
    plane_id INT NOT NULL REFERENCES plane(id),
    route_id INT NOT NULL REFERENCES route(id),
    departure_time TIMESTAMP,
    arrival_time TIMESTAMP,
    available_seats INT,
    transfer_flight_id INT REFERENCES flight(id)
);

CREATE TABLE seat(
    id SERIAL PRIMARY KEY,
    class VARCHAR(30),
    seat_number VARCHAR(5),
    user_id INT REFERENCES users(id),
    flight_id INT NOT NULL REFERENCES flight(id)
);

CREATE TABLE ticket(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    flight_id INT NOT NULL REFERENCES route(id),
    price DECIMAL,
    seat_id INT NOT NULL REFERENCES seat(id)    
);

CREATE TABLE user_tickets(
	user_id INT NOT NULL REFERENCES users(id),
    ticket_id INT NOT NULL REFERENCES ticket(id)
);

END;