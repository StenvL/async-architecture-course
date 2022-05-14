CREATE TABLE tasks (
    id                          SERIAL PRIMARY KEY,
    public_id                   UUID NOT NULL,
    title                       VARCHAR(100) NOT NULL,
    status                      VARCHAR(50) NOT NULL DEFAULT 'new',
    created                     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description                 TEXT NULL,
    assignee                    INT REFERENCES users
);