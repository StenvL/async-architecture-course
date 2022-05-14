CREATE TABLE users
(
    id    INT PRIMARY KEY,
    email VARCHAR(50)  NOT NULL,
    name  VARCHAR(100) NULL,
    role  VARCHAR(50)  NOT NULL
);

CREATE TABLE tasks
(
    id          UUID PRIMARY KEY,
    created     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    assignee    INT REFERENCES users NOT NULL,
    status      VARCHAR(50)          NOT NULL,
    title       VARCHAR(100)         NOT NULL,
    description TEXT                 NULL,
    cost        DECIMAL(8, 2)        NOT NULL,
    reward      DECIMAL(8, 2)        NOT NULL
);

CREATE TABLE accounts
(
    id      SERIAL PRIMARY KEY,
    user_id INT REFERENCES users NOT NULL,
    balance DECIMAL(10, 2)       NOT NULL
);

CREATE TABLE payments
(
    id         SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts NOT NULL,
    amount     DECIMAL(10, 2)          NOT NULL
);

CREATE TABLE accounts_audit_log
(
    id             SERIAL PRIMARY KEY,
    account_id     INT REFERENCES accounts NOT NULL,
    task_id        UUID REFERENCES tasks   NULL,
    type           VARCHAR(50)             NOT NULL,
    balance_change DECIMAL(10, 2)          NOT NULL,
    created        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
