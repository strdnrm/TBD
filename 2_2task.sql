--1
SELECT st.name, st.surname, hb.name
FROM student AS st
JOIN student_hobby AS sh ON st.id = sh.student_id
JOIN hobby AS hb ON hb.id = sh.hobby_id;
--2
SELECT *,
CASE
	WHEN sh.finished_at IS NOT NULL THEN sh.finished_at - sh.started_at
	ELSE NOW() - sh.started_at
END
AS Duration
FROM student AS st
LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
WHERE sh.id IS NOT NULL
ORDER BY duration DESC LIMIT 1;
--3
SELECT st.name, st.surname, st.id, st.age FROM student AS st
JOIN student_hobby AS sh ON sh.student_id = st.id 
JOIN hobby AS hb ON hb.id = sh.hobby_id
WHERE score > (SELECT AVG(score) FROM student) 
GROUP BY st.id HAVING SUM(hb.risk) > 0.9;
--4
SELECT st.name, st.surname, st.id, st.age, hb.name, sh.finished_at - sh.started_at
FROM student AS st
JOIN student_hobby AS sh ON sh.student_id = st.id
JOIN hobby AS hb ON hb.id = sh.hobby_id
WHERE sh.finished_at - sh.started_at IS NOT NULL;
--5
SELECT st.name, st.surname, st.id, st.age 
FROM student AS st
JOIN student_hobby AS sh ON sh.student_id = st.id
JOIN hobby AS hb ON hb.id = sh.hobby_id
GROUP BY st.id HAVING SUM(hb.id) > 1;
--6
SELECT st.n_group, AVG(st.score)
FROM student AS st
JOIN student_hobby AS sh ON sh.student_id = st.id
JOIN hobby AS hb ON hb.id = sh.hobby_id
WHERE st.id IN
(
	SELECT student_id FROM student_hobby
)
GROUP BY st.n_group;
--7
SELECT hb.name, hb.risk, NOW() - sh.started_at, st.id
FROM hobby AS hb 
JOIN student_hobby AS sh ON sh.hobby_id = hb.id
JOIN student AS st ON st.id = sh.student_id 
WHERE NOW() - sh.started_at IN
(
	SELECT MAX(NOW() - started_at)
	FROM student_hobby
	WHERE finished_at IS NULL
);
--8
SELECT hb.name, MAX(st.score)
FROM student_hobby AS sh
JOIN hobby AS hb ON hb.id = sh.hobby_id
JOIN student AS st ON st.id = sh.student_id
GROUP BY hb.name
ORDER BY MAX(st.score) DESC LIMIT 1;
--9
SELECT hb.name, st.score
FROM hobby AS hb
JOIN student_hobby AS sh ON sh.hobby_id = hb.id
JOIN student AS st ON st.id = sh.student_id
WHERE st.score BETWEEN 2.5 AND 3.5
AND st.n_group / 1000 = 2;
--10
 
--11

--12
SELECT st.n_group / 1000 AS course, COUNT(DISTINCT sh.hobby_id)
FROM student AS st
LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
GROUP BY st.n_group / 1000
--13
SELECT st.id, st.surname, st.name, st.age, st.n_group / 1000 AS course
FROM student AS st
LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
WHERE sh.id IS NULL AND st.score >= 4.5
ORDER BY st.n_group / 1000, st.age DESC
--14
CREATE VIEW students AS
SELECT st.id, st.name, st.surname, st.n_group, NOW() - sh.started_at AS duration
FROM student AS st
LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
WHERE sh.id IS NOT NULL
AND sh.finished_at IS NULL
AND extract(year from age(now(), sh.started_at)) > 5
--15
SELECT hb.name, COUNT(sh.student_id)
FROM hobby AS hb
LEFT JOIN student_hobby AS sh ON sh.hobby_id = hb.id
GROUP BY hb.name
--16
SELECT hb.id 
FROM hobby AS hb
LEFT JOIN student_hobby AS sh ON sh.hobby_id = hb.id
GROUP BY hb.id
ORDER BY COUNT(sh.student_id) DESC LIMIT 1;
--17
SELECT *
FROM student AS st
LEFT JOIN student_hobby AS sh ON st.id = sh.student_id
WHERE sh.hobby_id IN
(
	SELECT hb.id 
	FROM hobby AS hb
	LEFT JOIN student_hobby AS sh ON sh.hobby_id = hb.id
	GROUP BY hb.id
	ORDER BY COUNT(sh.student_id) DESC LIMIT 1
)
--18
SELECT id
FROM hobby
ORDER BY risk DESC LIMIT 3;
--19
SELECT *,
CASE
	WHEN sh.finished_at IS NOT NULL THEN sh.finished_at - sh.started_at
	ELSE NOW() - sh.started_at
END
AS Duration
FROM student AS st
LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
WHERE sh.id IS NOT NULL
ORDER BY duration DESC LIMIT 10;
--20
SELECT DISTINCT n_group
FROM student 
WHERE n_group IN
(
	SELECT n_group
	FROM student AS st
	LEFT JOIN student_hobby AS sh ON sh.student_id = st.id
	WHERE sh.id IS NOT NULL
	ORDER BY 
	CASE
		WHEN sh.finished_at IS NOT NULL THEN sh.finished_at - sh.started_at
		ELSE NOW() - sh.started_at
	END
	DESC LIMIT 10
)
--21
CREATE OR REPLACE VIEW stview AS
SELECT id, surname, name FROM student
ORDER BY score DESC;
--22