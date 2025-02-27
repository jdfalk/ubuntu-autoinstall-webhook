-- 004_add_phase_to_ipxe.sql
-- Add the "phase" column to ipxe_configuration table.
ALTER TABLE ipxe_configurations
ADD COLUMN IF NOT EXISTS phase TEXT NOT NULL DEFAULT 'install';
-- Add the "phase" column to ipxe_history table.
ALTER TABLE ipxe_history
ADD COLUMN IF NOT EXISTS phase TEXT NOT NULL DEFAULT 'install';
-- Optionally, if you need to update rows based on certain criteria,
-- you could run additional UPDATE statements here. For example:
--
UPDATE ipxe_configurations
SET phase = 'post-install'
WHERE config = '#!ipxe\nexit\n';
--
UPDATE ipxe_history
SET phase = 'post-install'
WHERE config = '#!ipxe\nexit\n';
