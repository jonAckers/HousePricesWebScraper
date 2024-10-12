-- +goose Up
CREATE TABLE properties (
    id                     UUID PRIMARY KEY,
    bedrooms               INT,
    bathrooms              INT,
    description            TEXT,
    address                VARCHAR(100),
    latitude               DOUBLE PRECISION,
    longitude              DOUBLE PRECISION,
    type                   VARCHAR(50),
    listing_update_reason  VARCHAR(20),
    listing_update_date    TIMESTAMP,
    price_amount           INT,
    price_currency_code    VARCHAR(5),
    estate_agent_telephone VARCHAR(20),
    estate_agent_name      VARCHAR(100),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS properties;
