-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

-- Create users_permissions table for many to many relationship
CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigserial NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigserial NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

-- Add first two permission into the permissions table
INSERT INTO permissions (code)
VALUES
('movies:read'),
('movies:write');
