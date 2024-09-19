CREATE TABLE request_logs (
                              id SERIAL PRIMARY KEY,                          -- Auto-incremented primary key
                              user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key to users
                              ad_id INT NOT NULL REFERENCES ad_details(id) ON DELETE CASCADE,  -- Foreign key to ad_details
                              timestamp TIMESTAMP NOT NULL                     -- The timestamp of the ad request
);

-- Indexes for optimizing queries and joins
CREATE INDEX idx_request_logs_user_id ON request_logs(user_id);   -- Index for faster joins with the users table
CREATE INDEX idx_request_logs_ad_id ON request_logs(ad_id);       -- Index for faster joins with the ad_details table
CREATE INDEX idx_request_logs_timestamp ON request_logs(timestamp); -- Index for queries based on time (e.g., daily reports)
