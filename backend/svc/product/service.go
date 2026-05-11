package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mattcarp12/mec/internal/models"
	"github.com/redis/go-redis/v9"
)

// ==========================================
// COMMAND SIDE (Writes to DB, Updates Cache)
// ==========================================

func CreateProduct(ctx context.Context, p *models.Product) error {
	// 1. Write to the primary Database (Postgres)
	err := CreateProductDB(ctx, p)
	if err != nil {
		return err
	}

	// 2. Immediately hydrate the Redis cache so subsequent reads are instant.
	if err = CreateProductCache(ctx, p); err != nil {
		log.Printf("Failed to cache product %s: %v\n", p.ID, err)
	}

	return nil
}

// ==========================================
// QUERY SIDE (Reads from Cache, Fallback to DB)
// ==========================================

func GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	// 1. Attempt to read from the Query Store (Redis)
	cachedData, err := GetProductByIdCache(ctx, id)
	if err == nil {
		return cachedData, nil
	} else if err != redis.Nil {
		// An actual Redis error occurred (not just a "not found")
		// Log it, but don't fail. Fallback to DB.
		fmt.Printf("Redis error: %v\n", err)
	}

	// 2. CACHE MISS: Fallback to the Command Store (Postgres)
	product, err := GetProductByIdDB(ctx, id)
	if err != nil {
		return nil, err // Product doesn't exist
	}

	// 3. Hydrate the cache for the next user
	if err = CreateProductCache(ctx, product); err != nil {
		log.Printf("Failed to cache product %s: %v\n", id, err)
	}

	return product, nil
}

