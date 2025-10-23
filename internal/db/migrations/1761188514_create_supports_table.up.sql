CREATE TABLE supports (
    id BIGSERIAL PRIMARY KEY,
    fan_id BIGINT NOT NULL,
    creator_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    FOREIGN KEY (fan_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT positive_support_amount CHECK (amount > 0),
    CONSTRAINT different_users CHECK (fan_id != creator_id)
);

CREATE INDEX idx_supports_fan_id ON supports(fan_id);
CREATE INDEX idx_supports_creator_id ON supports(creator_id);
CREATE INDEX idx_supports_status ON supports(status);