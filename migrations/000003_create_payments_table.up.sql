CREATE TABLE payments (
    id UUID PRIMARY KEY,
    
    -- A foreign key linking this payment directly to an order.
    -- 'ON DELETE CASCADE' means if an order is deleted, its payment record is also deleted.
    -- 'UNIQUE' ensures that an order can only have one payment record.
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    
    -- The amount to be paid.
    amount NUMERIC(10, 2) NOT NULL,
    
    -- The status of the payment attempt.
    status VARCHAR(50) NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- An index on order_id is automatically created because it's a UNIQUE key,
-- but creating one explicitly doesn't hurt. It's useful for finding a payment by order.
CREATE INDEX idx_payments_order_id ON payments (order_id);

-- Apply the timestamp-updating trigger to the 'payments' table.
CREATE TRIGGER update_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
