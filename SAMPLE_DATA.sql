-- Insert into Members Table
INSERT INTO Members (first_name, last_name, email, password_hash, date_of_birth, join_date, membership_type, status) VALUES 
('John', 'Doe', 'john.doe@example.com', 'hashed_password_1', '1990-01-15', NOW(), 'annual', 'active'),
('Jane', 'Smith', 'jane.smith@example.com', 'hashed_password_2', '1985-05-30', NOW(), 'monthly', 'active'),
('Alice', 'Johnson', 'alice.johnson@example.com', 'hashed_password_3', '1992-12-10', NOW(), 'lifetime', 'suspended'),
('Bob', 'Brown', 'bob.brown@example.com', 'hashed_password_4', '1988-04-25', NOW(), 'monthly', 'cancelled');

-- Insert into Membership Types Table
INSERT INTO MembershipTypes (type_name, duration, fee, benefits) VALUES 
('annual', 12, 120.00, 'Unlimited access for one year'),
('monthly', 1, 12.50, 'Unlimited access for one month'),
('lifetime', NULL, 1000.00, 'Unlimited access for lifetime');

-- Insert into Subscriptions Table
-- Assuming Member IDs and Membership Type IDs are generated as 1,2,3... in previous inserts
INSERT INTO Subscriptions (member_id, type_id, start_date, end_date, auto_renew) VALUES 
(1, 1, NOW(), NOW() + interval '1 year', TRUE),
(2, 2, NOW(), NOW() + interval '1 month', TRUE),
(3, 3, Now(), NULL, FALSE),
(4, 2, NOW() - interval '2 month', NOW() - interval '1 month', FALSE);

-- Insert into Payment Records Table
-- Assuming Payment Methods and Statuses are as defined
INSERT INTO PaymentRecords (member_id, amount, payment_date, payment_method, status) VALUES 
(1, 120.00, NOW(), 'credit card', 'completed'),
(2, 12.50, NOW() - interval '1 month', 'PayPal', 'completed'),
(3, 1000.00, NOW() - interval '5 years', 'bank transfer', 'completed'),
(4, 12.50, NOW(), 'credit card', 'failed');
