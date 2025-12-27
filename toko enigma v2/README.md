# Golang JWT

## Database

Buat database dan sesuaikan dengan file `config.go` untuk konfigurasinya, berikut untuk DDL nya:

```sql
CREATE DATABASE book_db;

CREATE TABLE mst_book
(
    id SERIAL primary key,
    title        VARCHAR(100),
    author       VARCHAR(100),
    release_year VARCHAR(4),
    pages        INTEGER
);

create table mst_user
(
    id SERIAL primary key,
    username VARCHAR(255) unique,
    password VARCHAR(255),
    role     VARCHAR(100)
);

INSERT INTO mst_user (username, password, role)
VALUES
    ('admin', 'admin', 'admin'),
    ('user', 'user', 'user');
```
