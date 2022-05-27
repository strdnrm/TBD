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

--Registration
INSERT INTO users (login, email, password, name, surname, is_admin)
VALUES ('login', 'email', 'password', 'name', 'surname', false);

--Route search by specified arrival date
SELECT * FROM flight
INNER JOIN route ON route.id = flight.id
WHERE flight.arrival_time::date = '2020.04.04'::date
AND route.source = 'MOSCOW'
AND route.destination = 'ISTANBUL';

--Route search by specified departure date
SELECT * FROM flight
INNER JOIN route ON route.id = flight.id
WHERE flight.departure_time::date = '2020.04.04'::date
AND route.source = 'MOSCOW'
AND route.destination = 'ISTANBUL';

--View available seats for the flight
SELECT * FROM seat
LEFT JOIN flight ON flight.id = seat.flight_id
LEFT JOIN route ON flight.route_id = route.id
WHERE user_id IS NULL;

--Ticket purchase
INSERT INTO ticket (user_id, flight_id, price, seat_id)
VALUES (1, 1, 100, 1);

INSERT INTO user_tickets(user_id, ticket_id)
VALUES(1, 1);

UPDATE seat
SET user_id = 1
WHERE seat_number = '12a';

--Adding new routes
INSERT INTO route (source,  destination)
VALUES ('NEW YORK', 'TOKYO');

--View user's flight statistics with departure date search
SELECT * 
FROM user_tickets
LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
LEFT JOIN flight ON flight.id = ticket.flight_id
LEFT JOIN route ON route.id = flight.route_id
WHERE user_tickets.user_id = 1
AND flight.departure_time::date = '2020.04.04'::date
ORDER BY flight.departure_time DESC;

--View user's flight statistics with arrival time date search
SELECT * 
FROM user_tickets
LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
LEFT JOIN flight ON flight.id = ticket.flight_id
LEFT JOIN route ON route.id = flight.route_id
WHERE user_tickets.user_id = 1
AND flight.arrival_time::date = '2020.04.04'::date
ORDER BY flight.departure_time DESC;

--View the user's flight statistics by departure point
SELECT * 
FROM user_tickets
LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
LEFT JOIN flight ON flight.id = ticket.flight_id
LEFT JOIN route ON route.id = flight.route_id
WHERE user_tickets.user_id = 1
AND route.source = 'PARIS'
ORDER BY flight.departure_time DESC;

--View the user's flight statistics by arrival point
SELECT * 
FROM user_tickets
LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
LEFT JOIN flight ON flight.id = ticket.flight_id
LEFT JOIN route ON route.id = flight.route_id
WHERE user_tickets.user_id = 1
AND route.destination = 'LONDON'
ORDER BY flight.departure_time DESC;

--Subscription to change the price of a given route
SELECT ticket.price, route.source, route.destination
FROM ticket
LEFT JOIN flight ON ticket.flight_id = flight.id
LEFT JOIN route ON flight.route_id = route.id
WHERE route.source = 'MOSCOW'
AND route.destination = 'TBILISI';

--Subscription for the appearance of tickets on the specified route on given dates
SELECT *
FROM flight
LEFT JOIN route ON flight.route_id = route.id 
WHERE  flight.departure_time::date = '2020.04.04'::date
AND route.source = 'ROME'
AND route.destination = 'BERLIN';