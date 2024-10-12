-- +goose Up
-- Step 1: Add a new column with INT type
ALTER TABLE properties
ADD COLUMN new_id INT;

-- Step 2: Drop the old UUID column
ALTER TABLE properties
DROP COLUMN id;

-- Step 3: Rename the new column to `id`
ALTER TABLE properties
RENAME COLUMN new_id TO id;

-- Step 4: Set the new `id` column as the primary key
ALTER TABLE properties
ADD PRIMARY KEY (id);

-- +goose Down
-- Step 1: Add the old `id` column with UUID type
ALTER TABLE properties
ADD COLUMN new_id UUID;

-- Step 2: Drop the `id` column
ALTER TABLE properties
DROP COLUMN id;

-- Step 3: Rename `new_id` back to id
ALTER TABLE properties
RENAME COLUMN new_id TO id;

-- Step 4: Re-add the primary key constraint on the id column
ALTER TABLE properties
ADD PRIMARY KEY (id);
