-- Members Table
CREATE TABLE Members
(
    member_id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    date_of_birth DATE,
    join_date TIMESTAMP(2) NOT NULL,
    membership_type VARCHAR(255),
    status VARCHAR(255),
    created_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Membership Types Table
CREATE TABLE MembershipTypes
(
    type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(255) UNIQUE NOT NULL,
    duration INT,
    fee NUMERIC(10, 2) NOT NULL,
    benefits TEXT,
    created_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Subscriptions Table
CREATE TABLE Subscriptions
(
    subscription_id SERIAL PRIMARY KEY,
    member_id INT NOT NULL REFERENCES Members(member_id),
    type_id INT NOT NULL REFERENCES MembershipTypes(type_id),
    start_date TIMESTAMP(2) NOT NULL,
    end_date TIMESTAMP,
    auto_renew BOOLEAN,
    created_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Payment Records Table
CREATE TABLE PaymentRecords
(
    payment_id SERIAL PRIMARY KEY,
    member_id INT NOT NULL REFERENCES Members(member_id),
    amount NUMERIC(10, 2) NOT NULL,
    payment_date TIMESTAMP(2) NOT NULL,
    payment_method VARCHAR(255),
    status VARCHAR(255),
    created_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(2) NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE OR REPLACE FUNCTION moddatetime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_timestamp_members
    BEFORE
UPDATE ON Members
    FOR EACH ROW
EXECUTE PROCEDURE moddatetime
(updated_at);

CREATE TRIGGER update_timestamp_membership_types
    BEFORE
UPDATE ON MembershipTypes
    FOR EACH ROW
EXECUTE PROCEDURE moddatetime
(updated_at);

CREATE TRIGGER update_timestamp_subscriptions
    BEFORE
UPDATE ON Subscriptions
    FOR EACH ROW
EXECUTE PROCEDURE moddatetime
(updated_at);

CREATE TRIGGER update_timestamp_payment_records
    BEFORE
UPDATE ON PaymentRecords
    FOR EACH ROW
EXECUTE PROCEDURE moddatetime
(updated_at);
