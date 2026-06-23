INSERT INTO users (user_name, user_level, email, password, nip, user_id) VALUES
('Admin', 'admin', 'admin@visualmesin.com', '$2a$10$YhYl4DoG.zX0glxrWcn7ceTQE0xTkHgBKRgVet0QtNoV9SfiUeAqO', 'ADM001', 'admin'),
('User Produksi', 'prod', 'user@visualmesin.com', '$2a$10$N4NWbjanR4ObNTikOqmo1uWZ3YK.VkA3LA0Cif9WQHBjLfzRs7oIW', 'PRD001', 'user')
ON DUPLICATE KEY UPDATE password = VALUES(password);
