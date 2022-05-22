CREATE TABLE users
(
    id    INT PRIMARY KEY,
    email VARCHAR(50)  NOT NULL,
    name  VARCHAR(100) NULL,
    role  VARCHAR(50)  NOT NULL
);

CREATE TABLE balance_changes
(
    account_id       INT,
    balance_changing DECIMAL(8, 2),
    "timestamp"      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks
(
    id          UUID PRIMARY KEY,
    reward      DECIMAL(8, 2),
    "timestamp" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);