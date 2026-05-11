package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mattcarp12/mec/internal/db"
	"github.com/mattcarp12/mec/internal/models"
)

func CreateProductDB(ctx context.Context, product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price_cents, inventory_count, attributes)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`

	err := db.DB.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.PriceCents,
		product.InventoryCount,
		product.Attributes,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func CreateProductCache(ctx context.Context, product *models.Product) error {
	// Cache the product in Redis with a TTL of 1 hour
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return db.Redis.Set(ctx, "product:"+product.ID.String(), data, 1*time.Hour).Err()
}

func GetProductByIdDB(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	query := `SELECT id, name, description, price_cents, inventory_count, attributes, created_at, updated_at FROM products WHERE id = $1`

	err := db.DB.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, err // Will return sql.ErrNoRows if not found
	}
	return &product, nil
}

func GetProductByIdCache(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	data, err := db.Redis.Get(ctx, "product:"+id.String()).Result()
	if err != nil {
		return nil, err // Will return redis.Nil if not found
	}
	err = json.Unmarshal([]byte(data), &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetAllProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := db.DB.SelectContext(ctx, &products, "SELECT * FROM products")
	return products, err
}
