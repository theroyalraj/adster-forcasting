CREATE TABLE users (
                       id SERIAL PRIMARY KEY,                           -- Auto-incremented primary key
                       user_id VARCHAR(100) NOT NULL,                   -- Unique user identifier
                       ip VARCHAR(45) NOT NULL,                         -- IP address of the user
                       geo_country VARCHAR(10) NOT NULL,                -- Country code (e.g., US)
                       geo_region VARCHAR(50),                          -- Region or state (e.g., CA)
                       geo_city VARCHAR(100),                           -- City (e.g., San Francisco)
                       device_type INT,                                 -- Device type (1: Mobile, 2: Tablet, 3: Desktop)
                       os VARCHAR(50),                                  -- Operating system (e.g., iOS, Android)
                       browser VARCHAR(50),                             -- Browser (e.g., Chrome, Safari)

    -- Ensuring uniqueness for the user_id
                       UNIQUE (user_id)
);

-- Indexes for optimizing common queries
CREATE INDEX idx_users_geo_country ON users(geo_country);     -- For filtering by country
CREATE INDEX idx_users_device_type ON users(device_type);     -- For filtering by device type
CREATE INDEX idx_users_ip ON users(ip);                       -- For filtering by IP address
