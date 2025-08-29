-- Create database
CREATE DATABASE tshirt_db;

-- Create user with password
CREATE USER tshirt WITH PASSWORD 'tshirt';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE tshirt_db TO tshirt;

-- Connect to the database
\c tshirt_db

-- Grant schema privileges
GRANT USAGE ON SCHEMA public TO tshirt;
GRANT CREATE ON SCHEMA public TO tshirt;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tshirt;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO tshirt;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tshirt;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tshirt;