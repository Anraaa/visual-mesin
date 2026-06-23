INSERT INTO users (user_name, user_level, email, password, nip, user_id) VALUES
('Admin', 'admin', 'admin@visualmesin.com', '$2a$10$RbCrEw5ug959rM0dcBSeO.wrd3Y9ZwEwEvZTZEMM3UA3Rbiaji4hK', 'ADM001', 'admin'),
('User Produksi', 'prod', 'user@visualmesin.com', '$2a$10$IM5kUnc8WfUmE6gJuwprquZIMYw2U5Vq8wz.2Am6/EB6PYI0FwmC2', 'PRD001', 'user')
ON DUPLICATE KEY UPDATE password = VALUES(password);
