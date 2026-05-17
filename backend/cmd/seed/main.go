package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mattcarp12/mec/internal/platform/config"
	"github.com/mattcarp12/mec/internal/platform/db"
	"github.com/sqlc-dev/pqtype"
)

func main() {
	ctx := context.Background()

	// 1. Load Config & Connect
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	dbConn, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// 2. Create an Organization
	org, err := queries.CreateOrganization(ctx, "Amorita Solutions LLC")
	if err != nil {
		log.Fatalf("Failed to create org: %v", err)
	}
	log.Printf("Created Organization: %s (ID: %v)", org.Name, org.ID)

	// 3. Create a System Admin User
	// Note: We use a dummy string for ID here. Later, this will be the Zitadel Subject ID.
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		ID:    "user_123_dummy_zitadel_id",
		Email: "admin@amoritasolutions.com",
		// OrgID: org.ID,
		OrgID: uuid.NullUUID{org.ID, true},
		Role:  "sys_admin",
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created User: %s", user.Email)

	// 4. Create Products
	products := []db.CreateProductParams{
		{
			OrgID:          org.ID,
			Name:           "Meridian Trail Running Shoes",
			Description:    "High-performance footwear for outdoor spaces and local parks.",
			PriceCents:     12500, // $125.00
			InventoryCount: 50,
			Attributes: pqtype.NullRawMessage{
				RawMessage: []byte(`{"color": "Slate", "sizes": ["10", "11", "12"]}`),
			},
		},
		{
			OrgID:          org.ID,
			Name:           "Air Fryer Cookbook",
			Description:    "100 slider bun and whole wheat recipes for the modern kitchen.",
			PriceCents:     2499, // $24.99
			InventoryCount: 200,
			Attributes: pqtype.NullRawMessage{
				RawMessage: []byte(`{"format": "Hardcover", "pages": 150}`),
			},
		},
	}

	for _, p := range products {
		prod, err := queries.CreateProduct(ctx, p)
		if err != nil {
			log.Printf("Failed to create product %s: %v", p.Name, err)
			continue
		}
		log.Printf("Created Product: %s - $%d (Cents)", prod.Name, prod.PriceCents)
	}

	log.Println("Database seeding completed successfully.")
}
