-- name: CreateProperty :one
INSERT INTO properties (
    id,
    bedrooms,
    bathrooms,
    description,
    address,
    latitude,
    longitude,
    type,
    listing_update_reason,
    listing_update_date,
    price_amount,
    price_currency_code,
    estate_agent_telephone,
    estate_agent_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;
--

-- name: GetPropertyById :one
SELECT * FROM properties
WHERE id=$1;
--
