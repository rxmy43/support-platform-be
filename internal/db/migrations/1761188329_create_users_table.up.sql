CREATE TYPE user_role AS ENUM ('fan', 'creator');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'fan'
);