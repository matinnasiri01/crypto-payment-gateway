-- ============================================================
-- 000001_init.down.sql
-- ============================================================

DROP INDEX IF EXISTS idx_invoices_expired_at;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_user_id;

DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS users;