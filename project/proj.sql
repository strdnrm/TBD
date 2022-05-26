CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	login VARCHAR(128),
	email VARCHAR(128),
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
    destintaion VARCHAR(128)
);

CREATE TABLE flight(
    id SERIAL PRIMARY KEY,
    plane_id INT NOT NULL REFERENCES plane(id),
    route_id INT NOT NULL REFERENCES route(id),
    departure_time TIMESTAMP,
    arrival_time TIMESTAMP,
    available_seats INT
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
    route_id INT NOT NULL REFERENCES route(id),
    price DECIMAL,
    seat_id INT NOT NULL REFERENCES seat(id)    
);

CREATE TABLE user_tickets(
	user_id INT NOT NULL REFERENCES users(id),
    ticket_id INT NOT NULL REFERENCES ticket(id)
);

--Registration
INSERT INTO users (login, email, password, name, surname, is_admin)
VALUES ('login', 'email', 'password', 'name', 'surname', false)

SELECT * FROM flight
WHERE arrival_time::date = '2020.04.04'::date;

SELECT * FROM flight
WHERE departure_time::date = '2020.04.04'::date;

--
SELECT * FROM seat
LEFT JOIN flight ON flight.id = seat.flight_id
WHERE user_id IS NULL;

INSERT INTO user_tickets(user_id, ticket_id)
VALUES(1, 1);