-- Drop trigger first
DROP TRIGGER IF EXISTS trigger_op_setup_updated_at ON operation_setup;

-- Drop function
DROP FUNCTION IF EXISTS update_op_setup_updated_at;

-- Drop table
DROP TABLE IF EXISTS operation_setup;