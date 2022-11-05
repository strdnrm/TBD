--single table sql queries
--1
SELECT name, surname FROM student WHERE score BETWEEN 4 AND 4.5;
SELECT name, surname FROM student WHERE score <= 4.5 AND score >= 4; 
--2
SELECT * FROM student WHERE n_group::VARCHAR LIKE '2%';
--3
SELECT * FROM student ORDER BY n_group DESC, name;
--4
SELECT * FROM student WHERE score > 4 ORDER BY score DESC;
--5
SELECT name, risk FROM hobby WHERE name LIKE 'Футбол' OR name LIKE 'Хоккей';
--6
SELECT hobby_id, student_id FROM student_hobby
WHERE (started_at BETWEEN '2017-01-01' AND '2019-01-01')
AND finished_at IS NOT NULL;
--7
SELECT * FROM student WHERE score >= 4.5 ORDER BY score DESC;
--8
SELECT * FROM student WHERE score >= 4.5 ORDER BY score DESC LIMIT 5;
--9
SELECT name,
CASE
	WHEN risk < 2 THEN 'Очень низкий'
	WHEN risk >= 2 AND risk < 4 THEN 'Низкий'
	WHEN risk >= 4 AND risk < 6 THEN 'Средний'
	WHEN risk >= 6 AND risk < 8 THEN 'Высокий'
	WHEN risk >= 8 THEN 'Очень высокий'
END
FROM hobby;
--10
SELECT * FROM hobby ORDER BY risk DESC LIMIT 3;