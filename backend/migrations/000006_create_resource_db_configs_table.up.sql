CREATE TABLE resource_db_configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    resource_name VARCHAR(255) NOT NULL UNIQUE,
    label VARCHAR(255),
    driver VARCHAR(20) NOT NULL DEFAULT 'mariadb',
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL DEFAULT 3306,
    database_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_last_test_success BOOLEAN,
    last_tested_at TIMESTAMP NULL,
    last_test_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active),
    INDEX idx_resource_name (resource_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
