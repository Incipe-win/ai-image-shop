-- Create database
CREATE DATABASE creative_studio_db;

-- Create user with password
CREATE USER creative WITH PASSWORD 'creative';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE creative_studio_db TO creative;

-- Connect to the database
\c creative_studio_db

-- Grant schema privileges
GRANT USAGE ON SCHEMA public TO creative;
GRANT CREATE ON SCHEMA public TO creative;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO creative;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO creative;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO creative;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO creative;