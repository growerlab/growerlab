-- /* create database */
-- CREATE DATABASE `growerlab` DEFAULT CHARACTER SET = `utf8mb4`;

-- /* create database user */
-- CREATE USER 'growerlab'@'localhost' IDENTIFIED BY 'growerlab';
-- ALTER USER 'growerlab'@'localhost' WITH MAX_QUERIES_PER_HOUR 0 MAX_UPDATES_PER_HOUR 0 MAX_CONNECTIONS_PER_HOUR 0;
-- GRANT DELETE, SELECT, EXECUTE, CREATE ROUTINE, ALTER ROUTINE, GRANT OPTION, REFERENCES, CREATE VIEW, TRIGGER, UPDATE, DROP, CREATE, LOCK TABLES, EVENT, INDEX, ALTER, SHOW VIEW, INSERT, CREATE TEMPORARY TABLES ON `growerlab`.* TO 'growerlab'@'localhost';


/* default password growerlab */
CREATE ROLE growerlab PASSWORD 'SCRAM-SHA-256$4096:lqvHjGo2ko6DWGB60ADcZg==$ZVQrnAK2nOGk0yL17+fKchLX+pK/78C3yYWVaDztGik=:UGRd+HjIww83HJ6rhkYSeWfRTumCx5Ttj8bSUvhcBaM=' INHERIT LOGIN;

/* create database for growerlab */
CREATE DATABASE growerlab OWNER growerlab ENCODING 'UTF8';