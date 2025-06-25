CREATE TABLE orders (
    id UUID PRIMARY KEY,
    
    -- The ID of the user who placed the order. Not a foreign key in this design,
    -- as user management could be a separate service.
    user_id UUID NOT NULL,
    
    -- The list of items in the order, including product ID, quantity, and price at the time of purchase.
    -- Storing this as JSONB is flexible and powerful.
    items JSONB NOT NULL,
    
    -- The total calculated price of the order.
    total_price NUMERIC(10, 2) NOT NULL,
    
    -- The current status of the order lifecycle.
    status VARCHAR(50) NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on user_id to quickly find all orders for a specific user.
CREATE INDEX idx_orders_user_id ON orders (user_id);

-- Index on status to quickly find all orders with a certain status (e.g., all 'PENDING' orders).
CREATE INDEX idx_orders_status ON orders (status);

-- Apply the same timestamp-updating trigger to the 'orders' table.
CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
