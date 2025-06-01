-- Group or category of items (e.g., "Essentials", "Group 1", "Green Products")
CREATE TABLE user_grocery_products.item_groups (
 id SERIAL PRIMARY KEY,
 name TEXT NOT NULL,
 slug TEXT UNIQUE, -- for filtering in frontend (e.g., "group-1")
 display_order INTEGER DEFAULT 0 -- to control display sequence in UI
);

-- Items that belong to groups
CREATE TABLE user_grocery_products.items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT,
    rating NUMERIC(2,1) CHECK (rating BETWEEN 0 AND 5), -- e.g., 4.5 stars
    group_id INTEGER REFERENCES user_grocery_products.item_groups(id) ON DELETE CASCADE,
    price NUMERIC(10, 2),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT,           -- New: detailed product description
    sustainability_facts TEXT   -- New: e.g., "Plastic-free packaging, locally sourced"
    ----- we'll need vendor_ec2 references like quantity and location as well
);

-- Dashboard banner/poster
CREATE TABLE user_grocery_products.dashboard_posters (
   id SERIAL PRIMARY KEY,
   image_url TEXT NOT NULL,
   link_url TEXT,
   active BOOLEAN DEFAULT TRUE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);