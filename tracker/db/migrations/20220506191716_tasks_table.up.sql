CREATE TABLE tasks (
    id                          SERIAL PRIMARY KEY,
    title                       VARCHAR(100) NOT NULL,
    status                      VARCHAR(50) NOT NULL DEFAULT 'new',
    craated                     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description                 TEXT NULL,
    assignee                    INT REFERENCES users
);