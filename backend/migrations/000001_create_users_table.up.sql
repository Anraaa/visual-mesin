CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nip VARCHAR(50) UNIQUE,
    user_id VARCHAR(100) UNIQUE,
    user_name VARCHAR(100) NOT NULL,
    user_level ENUM('admin','eng','tech','prod') NOT NULL DEFAULT 'prod',
    email VARCHAR(255) UNIQUE,
    email_verified_at TIMESTAMP NULL,
    password VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    remember_token VARCHAR(100),
    department VARCHAR(100),
    jabatan VARCHAR(100),
    themes_settings JSON,
    timestamp TIMESTAMP(3) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
