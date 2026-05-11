CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    -- E-commerce Rule: NEVER use floats for money. Store cents as integers.
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0), 
    inventory_count INTEGER NOT NULL CHECK (inventory_count >= 0),

    attributes JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast searches by name
CREATE INDEX idx_products_name ON products(name);
-- This allows you to query WHERE attributes ->> 'color' = 'red' instantly.
CREATE INDEX idx_products_attributes ON products USING GIN (attributes);