CREATE TABLE users (
    id             INT PRIMARY KEY,
    email          VARCHAR(50) NOT NULL,
    name           VARCHAR(100) NULL,
    role           VARCHAR(50) NOT NULL
);