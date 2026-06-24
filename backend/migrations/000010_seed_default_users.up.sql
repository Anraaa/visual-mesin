INSERT INTO users (user_name, user_level, email, password, nip, user_id) VALUES
('Admin', 'admin', 'admin@admin.com', '$2a$10$9USQDx35J2BpL6lsGr5NqOMOy5O1VnneR131Y4SQzUkoTfW0BkLtG', 'm26-134', 'admin'),
('User Produksi', 'prod', 'user@visualmesin.com', '$2a$10$IM5kUnc8WfUmE6gJuwprquZIMYw2U5Vq8wz.2Am6/EB6PYI0FwmC2', 'PRD001', 'user')
ON DUPLICATE KEY UPDATE password = VALUES(password);
