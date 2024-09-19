CREATE TABLE ad_details (
                            id SERIAL PRIMARY KEY,                           -- Auto-incremented primary key
                            domain VARCHAR(255) NOT NULL,                    -- Domain where the ad was shown (e.g., example.com)
                            url TEXT NOT NULL,                               -- Full URL of the page showing the ad
                            ad_position VARCHAR(50),                         -- Ad position (e.g., ATF: Above the Fold, BTF: Below the Fold)
                            ad_size VARCHAR(50),                             -- Size of the ad (e.g., 300x250, 728x90)

    -- Ensure uniqueness across domain, url, ad_position, and ad_size to avoid duplicate entries
                            UNIQUE (domain, url, ad_position, ad_size)
);

-- Indexes for optimizing common queries
CREATE INDEX idx_ad_details_domain ON ad_details(domain);      -- For filtering by domain
CREATE INDEX idx_ad_details_url ON ad_details(url);            -- For filtering by URL
