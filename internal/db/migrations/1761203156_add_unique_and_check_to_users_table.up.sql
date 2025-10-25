ALTER TABLE users
ADD CONSTRAINT users_phone_unique UNIQUE (phone);

ALTER TABLE users
ADD CONSTRAINT users_role_check CHECK (role IN ('fan', 'creator'));