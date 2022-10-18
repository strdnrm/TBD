--1
CREATE OR REPLACE FUNCTION func() RETURNS int AS $$
BEGIN
    RAISE NOTICE 'MESSAGE';
END;
$$ LANGUAGE plpgsql;

select func();
--2
CREATE OR REPLACE FUNCTION func() RETURNS date AS $$
DECLARE
	dt date;
BEGIN
    SELECT CURRENT_DATE INTO dt;
	RETURN dt;
END;
$$ LANGUAGE plpgsql;

select func();
--3
CREATE OR REPLACE FUNCTION plus(a int, b int) RETURNS int AS $$
DECLARE
	res int;
BEGIN
    SELECT a + b INTO res;
	RETURN res;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION multiply(a int, b int) RETURNS int AS $$
DECLARE
	res int;
BEGIN
    SELECT a * b INTO res;
	RETURN res;
END;
$$ LANGUAGE plpgsql;

select plus(3, 5);
select multiply(3, 5);

--4
CREATE OR REPLACE PROCEDURE mark(IN a int)
LANGUAGE plpgsql
AS $$
BEGIN
    IF a = 5 THEN RAISE NOTICE 'Отлично';
  	ELSIF a = 4 THEN RAISE NOTICE 'Хорошо';
  	ELSIF a = 3 THEN RAISE NOTICE 'Удовлетворительно';
	ELSIF a = 2 THEN RAISE NOTICE 'Неуд';
  	ELSE RAISE NOTICE 'Введенная оценка не верна';
  	END IF;
END;
$$ 

CREATE OR REPLACE PROCEDURE mark_w(IN a int)
LANGUAGE plpgsql
AS $$
BEGIN
	CASE a
		WHEN 5 THEN RAISE NOTICE 'Отлично';
		WHEN 4 THEN RAISE NOTICE 'Хорошо';
		WHEN 3 THEN RAISE NOTICE 'Удовлетворительно';
		WHEN 2 THEN RAISE NOTICE 'Неуд';
		ELSE RAISE NOTICE 'Введенная оценка не верна';
  	END CASE;
END;
$$ 

call mark(6);
call mark_w(2);

--5
DO $$
DECLARE 
	a int := 20;
BEGIN
	LOOP
		IF a > 30 THEN EXIT;
		END IF;
		RAISE NOTICE '%', a^2;
		a = a + 1;
	END LOOP;
END;
$$;

DO $$
DECLARE 
	a int := 20;
BEGIN
	WHILE a < 31 
	LOOP
		RAISE NOTICE '%', a^2;
		a = a + 1;
	END LOOP;
END;
$$;


DO $$
BEGIN
	FOR a IN 20..30 LOOP
		RAISE NOTICE '%', a^2;
	END LOOP;
END;
$$;
--6
CREATE OR REPLACE FUNCTION сollatz(a int) RETURNS int
AS $$
BEGIN
    WHILE a <> 1 LOOP
		IF a % 2 = 0 THEN
			a = a / 2;
			RAISE NOTICE '%', a;
		ELSE
			a = a * 3 + 1;
			RAISE NOTICE '%', a;
		END IF;
	END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE сollatzp(IN a int)
LANGUAGE plpgsql
AS $$
BEGIN
	WHILE a <> 1 LOOP
		IF a % 2 = 0 THEN
			a = a / 2;
			RAISE NOTICE '%', a;
		ELSE
			a = a * 3 + 1;
			RAISE NOTICE '%', a;
		END IF;
	END LOOP;
END;
$$ 

call сollatzp(6);
select * from сollatz(12);
--7
CREATE OR REPLACE FUNCTION lucas(n int) RETURNS int
AS $$
DECLARE
i int;
a int;
b int;
tmp int;
BEGIN
	b = 1;
	a = 2;
	IF n = 0 THEN
		RETURN 2;
	END IF;
	IF n = 1 THEN 
		RETURN 1;
	END IF;
	i = 2;
    WHILE i < n LOOP
		tmp := b;
		b = b + a;
		a = tmp;
		i = i + 1;
	END LOOP;
	RETURN b;
END;
$$ LANGUAGE plpgsql;

select lucas(5);

CREATE OR REPLACE PROCEDURE lucasp(IN n int)
LANGUAGE plpgsql
AS $$
DECLARE
i int;
a int;
b int;
tmp int;
BEGIN
	b = 1;
	a = 2;
	IF n = 0 THEN
		RAISE NOTICE '%|%', a, b;
		RETURN;
	END IF;
	IF n = 1 THEN 
		RAISE NOTICE '%|%', a, b;
		RETURN;
	END IF;
	i = 2;
    WHILE i < n LOOP
		tmp := b;
		b = b + a;
		a = tmp;
		i = i + 1;
		RAISE NOTICE '%|%', a, b;
	END LOOP;
END;
$$ 

call lucasp(5);
--8
CREATE OR REPLACE FUNCTION count_by_year(year int) RETURNS int
AS $$
DECLARE
peoplenum int;
BEGIN
	SELECT COUNT(id) INTO peoplenum
	FROM people
	WHERE EXTRACT(YEAR FROM people.birth_date) = year;
	RETURN peoplenum;
END
$$
LANGUAGE plpgsql;

select count_by_year(2003);
--9
CREATE OR REPLACE FUNCTION count_by_eyes(eye_color varchar) RETURNS int
AS $$
DECLARE
peoplenum int;
BEGIN
	SELECT COUNT(id) INTO peoplenum
	FROM people
	WHERE people.eyes = eye_color;
	RETURN peoplenum;
END
$$
LANGUAGE plpgsql;

select count_by_eyes('blue');
--10
AS $$
DECLARE
idy int;
BEGIN
	SELECT id INTO idy
	FROM people
	ORDER BY people.birth_date DESC LIMIT 1;
	RETURN idy;
END
$$
LANGUAGE plpgsql;

select get_the_youngest();
--11
CREATE OR REPLACE PROCEDURE body_mass_index(IN ind int)
LANGUAGE plpgsql
AS $$
DECLARE 
	r people%rowtype;
BEGIN
	FOR r IN 
		SELECT * FROM people
		WHERE people.weight / people.growth^2 > ind
	LOOP
		RAISE NOTICE 'ID: % NAME: % SURNAME: %', r.id, r.name, r.surname;
	END LOOP;
END
$$


call body_mass_index(0);
--12
BEGIN;
CREATE TABLE relatives (person_id int REFERENCES people(id), relative_id int REFERENCES people(id));
COMMIT;
END;
--13
CREATE OR REPLACE PROCEDURE new_person(IN new_name varchar, new_surname varchar,
											  new_birth_date date, new_growth real,
											  new_weight real, new_eyes varchar,
											  new_hair varchar, relative_id int)
LANGUAGE plpgsql
AS $$
DECLARE
	rid int;
BEGIN
	INSERT INTO people (name, surname, birth_date, growth, weight, eyes, hair)
	VALUES (new_name, new_surname, new_birth_date, new_growth, new_weight, new_eyes, new_hair)
	RETURNING id INTO rid;
 	INSERT INTO relatives(person_id, relative_id)
 	VALUES (rid, relative_id);
	INSERT INTO relatives(person_id, relative_id)
 	VALUES (relative_id, rid);
END;
$$;

call new_person('Lev', 'Muril', '1988-12-12'::date, 180.2, 70.4, 'green', 'blond', 1);
--14
BEGIN;
ALTER TABLE people ADD COLUMN update_date date;
COMMIT;
END;
--15
CREATE OR REPLACE PROCEDURE update_data(IN idn int, new_groth real, new_weight real)
LANGUAGE plpgsql
AS $$
BEGIN
	UPDATE people SET growth = new_groth, weight = new_weight, update_date = CURRENT_DATE
	WHERE people.id = idn;
END;
$$;
DROP PROCEDURE update_data(integer,real,real)

call update_data(2, 220, 50);