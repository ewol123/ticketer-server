INSERT INTO role(id,name)
VALUES
('72daf87a-fda4-4c72-aff9-85edd68d155f','user'),
('336a3ff6-9fdb-496f-ac8c-e37759969cf2','admin');

INSERT INTO "user"
(id,created_at,full_name,email,password,updated_at,status,registration_code,reset_password_code)
VALUES
('e66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'peti', 'peti@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null);

INSERT INTO "user_role"
(user_id,role_id)
VALUES
('e66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2')