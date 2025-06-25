DROP TABLE IF EXISTS products;
-- The trigger is dropped automatically when the table is dropped.
-- The function remains, but it's harmless. If you want to be perfectly clean, you can also add:
-- DROP FUNCTION IF EXISTS update_updated_at_column(); 
-- (But be careful if other tables use it).
