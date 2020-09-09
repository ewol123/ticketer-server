INSERT INTO role(id,name)
VALUES
('72daf87a-fda4-4c72-aff9-85edd68d155f','user'),
('336a3ff6-9fdb-496f-ac8c-e37759969cf2','admin');

INSERT INTO "user"
(id,created_at,full_name,email,password,updated_at,status,registration_code,reset_password_code)
VALUES
--GetUser test data
('a66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test2', 'test1@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
--ConfirmRegistration test data
('b66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test2', 'test2@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', '123456',null),
--DeleteUser test data
('d076f530-2453-4af2-a9a2-52b54dc3d36f', '2011-01-01 00:00:00', 'test3', 'test3@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
--GetAllUserTest test data
('d66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test4', 'test4@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
('e66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test5', 'test5@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
('f66c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test6', 'test6@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
--Login test data
('f16c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test7', 'test7@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
--ResetPassword test data
('f26c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test8', 'test8@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,'123456'),
--SendPasswdReset test data
('f36c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test9', 'test9@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null),
--UpdateUser test data
('f46c0b06-ec6c-45e1-8619-27e14c3ed92d', '2011-01-01 00:00:00', 'test10', 'test10@test.com', '$2y$10$3CaZ7aya5F0YJpWfuypnHOi.LNzCQP4lx33roaXZiKERBXRBOIJUG', '2011-01-01 00:00:00', 'active', null,null);

--dont forget the roles for tests
INSERT INTO "user_role" (user_id,role_id)
VALUES
('a66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('b66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('d076f530-2453-4af2-a9a2-52b54dc3d36f', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('d66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('e66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('f66c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('f16c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('f26c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('f36c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2'),
('f46c0b06-ec6c-45e1-8619-27e14c3ed92d', '336a3ff6-9fdb-496f-ac8c-e37759969cf2');