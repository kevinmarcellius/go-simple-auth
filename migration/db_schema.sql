-- Enable the uuid-ossp extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the users table
CREATE TABLE users (
    -- id: UUID primary key, automatically generated on creation
    id UUID PRIMARY KEY,

    -- username: User's chosen name, must be unique
    username VARCHAR(50) UNIQUE NOT NULL,

    -- email: User's email address, must be unique for login and communication
    email VARCHAR(255) UNIQUE NOT NULL,

    -- password_hash: Store a secure hash of the user's password, not the plain text
    password_hash VARCHAR(255) NOT NULL,

    -- is_admin: Flag to determine if the user has administrative privileges
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    -- created_at: Timestamp for when the user account was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- updated_at: Timestamp for the last time the user account was updated
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the email column for faster lookups during login
CREATE INDEX idx_users_email ON users(email);

-- Create a trigger function to automatically update the updated_at timestamp on any row modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Attach the trigger to the users table
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
