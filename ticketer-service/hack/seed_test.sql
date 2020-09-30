INSERT INTO "ticket"
(id,user_id,worker_id,fault_type,address,full_name,phone,geo_location,image_url,status,created_at,updated_at)
VALUES
--DeleteTicketAdmin test data
('d076f530-2453-4af2-a9a2-52b54dc3d36f','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00'),


--Create bunch of tickets
('f735dd08-bbdd-4a65-9336-df21804eb47e','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00'),
('fc21847d-2cec-403d-9ba4-c64e8756a400','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00'),
('101018d2-3ebf-4aab-b233-a4a5640822ae','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00'),
('e96e5e54-1c1f-4d46-bdcd-bfe346176c5c','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00'),
('cc773940-bb68-4a05-b70d-2df927c2bea8','180823c3-e9ea-4bb1-bdfb-a1cc1cdbf6bd', '60dd5185-6003-48da-9ff1-998a4477529c', 'leak', 'test address', 'Peter', '36300001111', ST_MakePoint('1.1','-1.1'), 'http://myimage/stg.jpg','draft', '2011-01-01 00:00:00', '2011-01-01 00:00:00')
