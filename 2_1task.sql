--group functions
--1
SELECT n_group, COUNT(n_group) FROM student GROUP BY n_group; 
--2
SELECT n_group, MAX(score) FROM student GROUP BY n_group; 
--3
SELECT surname, COUNT(surname) FROM student GROUP BY surname; 
--4
SELECT year_of_birth, count(id) FROM student GROUP BY year_of_birth;
--5
SELECT SUBSTR(n_group::VARCHAR, 1, 1) AS Курс, AVG(score) 
FROM student GROUP BY SUBSTR(n_group::VARCHAR, 1, 1);
--6
SELECT n_group, MAX(score) FROM student GROUP BY n_group ORDER BY MAX(score) LIMIT 1;
--7
SELECT n_group, AVG(score) FROM student GROUP BY n_group HAVING AVG(score) > 3.5 ORDER BY AVG(score) DESC;
--8
SELECT n_group, COUNT(n_group), MAX(score), AVG(score), MIN(score) FROM student GROUP BY n_group;
--9
SELECT * FROM student WHERE score IN (SELECT MAX(score) FROM student WHERE n_group = 2231 GROUP BY n_group);
--10
SELECT * FROM student WHERE (n_group, score) IN (SELECT n_group, MAX(score) FROM student GROUP BY n_group);