CREATE TABLE products (
    -- UUID for a unique product identifier.
    id UUID PRIMARY KEY,
    
    -- Name of the product, e.g., "Go Programming T-Shirt".
    name VARCHAR(255) NOT NULL,
    
    -- Price stored as a NUMERIC type to avoid floating-point inaccuracies.
    -- Allows for prices up to 99,999,999.99.
    price NUMERIC(10, 2) NOT NULL,
    
    -- The available stock quantity. 'DEFAULT 0' ensures it's never NULL.
    -- A CHECK constraint prevents it from ever going below zero.
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    
    -- Timestamps with timezone are best practice.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- This function automatically updates the 'updated_at' timestamp on any change.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply the trigger to the 'products' table.
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
