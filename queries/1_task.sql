

CREATE TABLE student(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	surname VARCHAR(255),
	age INT,
	n_group INT,
	course INT,
	score FLOAT
); 

CREATE TABLE hobby(
	id SERIAL PRIMARY KEY,
	name varchar(255) NOT NULL,
	risk NUMERIC(3, 2) NOT NULL
);

CREATE TABLE student_hobby(
	id SERIAL PRIMARY KEY,
	student_id INT NOT NULL REFERENCES student(id),
	hobby_id INT NOT NULL REFERENCES hobby(id),
	started_at DATE NOT NULL,
	finished_at TIMESTAMP
);

INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Гилмор', 'Пуф', 20, 2231, 2, 4.3);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Майкл', 'Карлеоне', 19, 2231, 2, 2.1);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Иван', 'Лоскутов', 23, 4412, 4, 3.2);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Виктория', 'Романова', 21, 3221, 3, 4.4);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Хорус', 'Луперкаль', 23, 4412, 4, 2.3);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Дмитрий', 'Великий', 18, 4412, 2, 4.3);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Мария', 'Шпиц', 22, 2231, 2, 4.9);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Владимир', 'Шлопков', 21, 2231, 2, 4.6);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Тимур', 'Смирный', 23, 4412, 4, 2.7);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Вильгелм', 'Строитель', 25, 3221, 3, 2.6);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Валерия', 'Деспотова', 18, 3221, 3, 4.1);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Семен', 'Мирный', 19, 2231, 2, 4.8);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Уильям', 'Рено', 17, 4412, 4, 3.9);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Виктория', 'Штольц', 20, 2231, 2, 2.6);
INSERT INTO student (name, surname, age, n_group, course, score) VALUES ('Геральт', 'Ривинов', 25, 4412, 4, 2.5);

INSERT INTO hobby (name, risk) VALUES ('Плавание', 0.4);
INSERT INTO hobby (name, risk) VALUES ('Прыжки с трамплина', 0.9);
INSERT INTO hobby (name, risk) VALUES ('Шахматы', 0.2);

INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (5, 1, '2015-12-21', '12-06-2018 10:12:33');
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (4, 1, '2017-11-15', '07-08-2021 23:32:12');
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (2, 2, '2012-02-29', null);
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (2, 3, '2012-03-04', '12-08-2013 11:22:33');
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (8, 1, '2020-10-08', null);
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (7, 2, '2017-07-14', null);
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (12, 3, '2018-05-23', null);
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (13, 1, '2018-01-31', '12-03-2019 12:55:44');
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (4, 3, '2019-03-21', null);
INSERT INTO student_hobby (student_id, hobby_id, started_at, finished_at) VALUES (4, 2, '2020-11-24', null);