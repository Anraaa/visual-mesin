-- 000015: Create orphan tables (bpbl, gtentire, rtl-tl1, recipe_material)

-- 4.1 bpbl (Master part/production problem tracking)
CREATE TABLE IF NOT EXISTS `bpbl` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `bc_entried` VARCHAR(50),
    `mcncode` VARCHAR(50),
    `probcode` VARCHAR(50),
    `jdge` VARCHAR(50),
    `opr` VARCHAR(50),
    `oprname` VARCHAR(50),
    `jdge_date` VARCHAR(50),
    `date_shift` VARCHAR(50),
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4.2 gtentire (Green tire tracking)
CREATE TABLE IF NOT EXISTS `gtentire` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `bc_entried` VARCHAR(50),
    `bld_date` VARCHAR(50),
    `bld_mcn01` VARCHAR(50),
    `bld_shift` VARCHAR(50),
    `cur_in` VARCHAR(50),
    `cur_out` VARCHAR(50),
    `cur_shift` VARCHAR(50),
    `cur_mcn` VARCHAR(50),
    `cur_opr01` VARCHAR(50),
    `jdge_date` VARCHAR(50),
    `jdge` VARCHAR(50),
    `probcode` VARCHAR(50),
    `pic` VARCHAR(50),
    `status` VARCHAR(50),
    `whs_in` VARCHAR(50),
    `whs_out` VARCHAR(50),
    `serialnumb` VARCHAR(50),
    `sarana` VARCHAR(50),
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4.3 rtl-tl1 (Production line speed monitoring)
CREATE TABLE IF NOT EXISTS `rtl-tl1` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `tgl` VARCHAR(50) NOT NULL DEFAULT '',
    `jam` VARCHAR(50) NOT NULL DEFAULT '',
    `shift` VARCHAR(10) NOT NULL DEFAULT '',
    `codesize` VARCHAR(50) NOT NULL DEFAULT '',
    `SpecSpeed` VARCHAR(50) NOT NULL DEFAULT '',
    `SpeedSP` VARCHAR(50) NOT NULL DEFAULT '',
    `SpeedSrv` VARCHAR(50) NOT NULL DEFAULT '',
    `SpeedProd` VARCHAR(50) NOT NULL DEFAULT '',
    `SpeedA` VARCHAR(50) NOT NULL DEFAULT '',
    `tol--` VARCHAR(20) NOT NULL DEFAULT '',
    `tol++` VARCHAR(20) NOT NULL DEFAULT '',
    `leader` VARCHAR(20) NOT NULL DEFAULT '',
    `operator1` VARCHAR(20) NOT NULL DEFAULT '',
    `operator2` VARCHAR(20) NOT NULL DEFAULT '',
    `operator3` VARCHAR(20) NOT NULL DEFAULT '',
    `operator4` VARCHAR(20) NOT NULL DEFAULT '',
    `operator5` VARCHAR(20) NOT NULL DEFAULT '',
    `average` VARCHAR(50) NOT NULL DEFAULT '',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4.4 recipe_material (Recipe ↔ Material junction)
CREATE TABLE IF NOT EXISTS `recipe_material` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `recipe_id` INT NOT NULL,
    `material_id` INT NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_recipe_material_recipe` (`recipe_id`),
    KEY `idx_recipe_material_material` (`material_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed resource_group_items for new tables
-- bpbl → Group 1 (Material)
INSERT IGNORE INTO resource_group_items (group_id, resource_name, label, sort_order) VALUES
    (1, 'bpbl', 'bpbl', 6);

-- gtentire → Group 2 (Curing)
INSERT IGNORE INTO resource_group_items (group_id, resource_name, label, sort_order) VALUES
    (2, 'gtentire', 'gtentire', 4);

-- rtl-tl1 → Group 6 (Trimming)
INSERT IGNORE INTO resource_group_items (group_id, resource_name, label, sort_order) VALUES
    (6, 'rtl-tl1', 'rtl-tl1', 3);

-- recipe_material → Group 4 (Recipe)
INSERT IGNORE INTO resource_group_items (group_id, resource_name, label, sort_order) VALUES
    (4, 'recipe_material', 'recipe_material', 4);
