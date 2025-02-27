-- 003_add_phase_to_cloud_init.sql

-- Add the "phase" column to cloud_init_user_data table with a default value.
ALTER TABLE cloud_init_userdata
ADD COLUMN phase TEXT NOT NULL DEFAULT 'install';

-- Add the "phase" column to cloud_init_history table with a default value.
ALTER TABLE cloud_init_history
ADD COLUMN phase TEXT NOT NULL DEFAULT 'install';

-- Optional: If you want to later differentiate between "install" and "post-install",
-- you may update rows manually or via further migration scripts.
