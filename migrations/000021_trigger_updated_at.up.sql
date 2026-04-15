-- =========================================
-- COMMON FUNCTION
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER FOR production_target
-- =========================================

DROP TRIGGER IF EXISTS trg_pt_updated_at ON production_target;

CREATE TRIGGER trg_pt_updated_at
BEFORE UPDATE ON production_target
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();