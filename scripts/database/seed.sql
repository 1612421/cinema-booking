-- seed.sql - Seed data for cinema booking system with UUIDs only

-- Toàn bộ user đều sử dụng password là test123
INSERT
IGNORE INTO users (id, username, password, status, address, phone_number, created_at, updated_at)
VALUES 
('0853150e-2426-4d94-a6ab-b41dc0303c78','nghiann03051','$2a$10$SzLQuR0ow7xZCz.z.r.w8OR.Z0mHyTDBqcF4vibQ6xCBbf7ky1vJa','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('3e32fb82-d37b-4e1a-9a63-f903c1c47195','nghiann03055','$2a$10$v349PHDwaA9qba1uLOEhFOh3hPh7anKwvxz2OHJXFW36P9AQPYv.u','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('512c9b1f-59b2-4047-8c64-032fcb0ec6de','nghiann030510','$2a$10$vGQA0nlW9ysZSjThfvVUCuQiQvKq5C76pugR1..egRczqYAS.r3Yu','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('524f113b-128b-4948-a73f-48e04b9a55d1','nghiann03056','$2a$10$6WPqMf6zjmPbHcIFztcUleHiMXsK5erQ3ZI359.eK413FGqLilYN6','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('75db16c0-cfae-4b3a-b14f-64e6c7b866d9','nghiann03052','$2a$10$pM.BcpTo58j0YBqnIYXyy.s4GEK/uDh4tg9j/namqYZeBks3oGVhy','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('8a449727-852f-46ba-a510-fb849f5f02e7','nghiann03058','$2a$10$mDWITXKYuV2PIHNQKYSOUuY1DOTcsgs2JOldidlVlGx7vRYj3ISuW','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('a1eda7ae-1251-401e-842b-d431612a7fe0','nghiann03059','$2a$10$FdkDob6urfk/wPfzXMAGn.VRicKAyZYMwcBvD/NYjDopAJ8Z7W/Di','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('a1eda7ae-1251-401e-842b-d431612a7fe1','user1','$2a$10$FdkDob6urfk/wPfzXMAGn.VRicKAyZYMwcBvD/NYjDopAJ8Z7W/Di','active','123 Test St','0123456789',NOW(),NOW()),
('a8e21cc6-d4fe-49f9-928f-1546122f90f2','nghiann03057','$2a$10$c5KqNTAbMXRf6d7A67309.fASLK.Eb5Asdsyc0if3istlum6deUm2','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('d182da50-3c56-4daf-bb3c-a2f050f0a48d','nghiann03053','$2a$10$ws.d.yJaEeSodyFk76plK.ZhpHwwAnMILSMwN.CtH1K02s1Dw6N7S','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('dbf8945d-037d-4630-a5e5-0521b0e6ee19','nghiann03054','$2a$10$NA9CEE1tUeVdsBOweSBGVenovYv/OccJGgWXEX4IfJOvkuWWvtfxS','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW()),
('f621a23e-1926-4675-96b2-2980dc1f78b3','nghiann0305','$2a$10$rCkefvbq6X2DItHTizhYoOb7jTEDLCPGgH3TXen7JOc4Djd4JBZlC','active','Cam Ranh, Khanh Hoa','0366773942',NOW(),NOW());

-- Insert test movie
INSERT
IGNORE INTO movies (id, title, duration, created_at, updated_at) VALUES
('8e5d1799-fb8d-4e92-a241-16ccf9c84aa0','Avengers: Endgame',181,NOW(),NOW()),
('b8edb4d3-a30a-4bdb-b3b6-150305c9fdef','Black Myth: Wukong',200,NOW(),NOW()),
('f15568f1-57d7-4ee0-895b-5a5fad83e941','Natra',200,NOW(),NOW());

-- Insert 500 seats for screen with UUID
CALL seed_seats('c8010e6b-e405-4e25-bc35-9006c6277a17', 10, 5);
CALL seed_seats('4f990f2e-3a8b-42e3-a3ef-42700d8f437c', 2, 5);

-- Insert showtime
INSERT
IGNORE INTO showtime (id, movie_id, screen_id, start_time, created_at, updated_at)
VALUES ('1c57ed08-8f95-4c8c-a63a-6e7d73db63e3', '8e5d1799-fb8d-4e92-a241-16ccf9c84aa0',
        'c8010e6b-e405-4e25-bc35-9006c6277a17', '2025-07-15 18:00:00', NOW(),
        NOW()),
       ('ccd13ac7-86e8-4f9e-96fd-5f41977eed6d', 'b8edb4d3-a30a-4bdb-b3b6-150305c9fdef',
        '4f990f2e-3a8b-42e3-a3ef-42700d8f437c', '2025-07-20 07:00:00', NOW(), NOW());

