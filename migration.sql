CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY,
    count INT NOT NULL
);

CREATE TABLE access_logs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) references users (username),
    ip_address VARCHAR(255) NOT NULL,
    data JSON DEFAULT '{}',
    access_time TIMESTAMPTZ DEFAULT NOW()
);