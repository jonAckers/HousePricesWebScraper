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

-- name: UpdatePropertyById :exec
UPDATE properties
SET
    bedrooms=$1,
    bathrooms=$2,
    description=$3,
    address=$4,
    latitude=$5,
    longitude=$6,
    type=$7,
    listing_update_reason=$8,
    listing_update_date=$9,
    price_amount=$10,
    price_currency_code=$11,
    estate_agent_telephone=$12,
    estate_agent_name=$13
WHERE
    id = $14
RETURNING *;
