-- name: CreateOrganization :one
INSERT INTO organizations (name) 
VALUES ($1) 
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (id, email, org_id, role) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (org_id, name, description, price_cents, inventory_count, attributes) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: ListProductsByOrg :many
SELECT * FROM products 
WHERE org_id = $1 
ORDER BY created_at DESC;

-- name: GetProductByID :one
SELECT * FROM products 
WHERE id = $1 LIMIT 1;