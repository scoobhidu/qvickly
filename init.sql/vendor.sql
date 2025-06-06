-- CREATE TYPE account_type AS ENUM ('store', 'restaurant');

-- CREATE TABLE accounts.accounts (
--   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   phone_number VARCHAR(15) NOT NULL UNIQUE,
--   account_type account_type NOT NULL,
--   business_name VARCHAR(100) NOT NULL,
--   owner_name VARCHAR(100),
--   email VARCHAR(100),
--   address TEXT,
--   latitude DOUBLE PRECISION,
--   longitude DOUBLE PRECISION,
--   gstin_number varchar(35),
--   opening_time TIME,
--   closing_time TIME,
--   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE accounts.restaurant_details (
--     account_id UUID PRIMARY KEY REFERENCES accounts.accounts(id) ON DELETE CASCADE,
--     fssai_license_no VARCHAR(50),
--     cuisine_id UUID NOT NULL REFERENCES constants.cuisines(id) ON DELETE CASCADE,
--     UNIQUE (account_id, cuisine_id)
-- );

-- CREATE TABLE constants.cuisines (
--   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   name VARCHAR(50) UNIQUE NOT NULL
-- );

-- CREATE TABLE vendor_ec2.items (
--    id SERIAL PRIMARY KEY,
--    account_id INTEGER NOT NULL REFERENCES vendor_ec2.accounts.accounts(id) ON DELETE CASCADE,
--    category_id INTEGER REFERENCES categories(id),
--    name TEXT NOT NULL,
--    description TEXT,
--    price_retail NUMERIC(10,2),
--    price_wholesale NUMERIC(10,2),
--    is_available BOOLEAN DEFAULT TRUE,
--    stock INTEGER DEFAULT 0,
--    created_at TIMESTAMP DEFAULT NOW(),
--    updated_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE items.categories (
--     id SERIAL PRIMARY KEY,
--     name TEXT UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE items.items (
--    id SERIAL PRIMARY KEY,
--    account_id uuid NOT NULL REFERENCES vendor_ec2.accounts.accounts(id) ON DELETE CASCADE,
--    category_id INTEGER REFERENCES items.categories(id),
--    name TEXT NOT NULL,
--    description TEXT,
--    price_retail NUMERIC(10,2),
--    price_wholesale NUMERIC(10,2),
--    is_available BOOLEAN DEFAULT TRUE,
--    stock INTEGER DEFAULT 0,
--    created_at TIMESTAMP DEFAULT NOW(),
--    updated_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE TABLE items.item_images (
--  id SERIAL PRIMARY KEY,
--  item_id INTEGER NOT NULL REFERENCES items.items(id) ON DELETE CASCADE,
--  image_url TEXT NOT NULL,
--  position INTEGER CHECK (position BETWEEN 1 AND 4),
--  created_at TIMESTAMP DEFAULT NOW()
-- );

-- CREATE INDEX idx_items_account_id ON items.items(account_id);
-- CREATE INDEX idx_items_category_id ON items.items(category_id);
-- CREATE INDEX idx_item_images_item_id ON items.item_images(item_id);

-- -- Delivery partners table
-- CREATE TABLE delivery_partners.delivery_partners (
--                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--                                    name VARCHAR(100) NOT NULL,
--                                    phone_number VARCHAR(15) NOT NULL,
--                                    pin VARCHAR(10),
--                                    is_active BOOLEAN DEFAULT true,
--                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
--
-- -- Add delivery assignment to orders
-- ALTER TABLE orders.orders ADD COLUMN delivery_partner_id UUID REFERENCES delivery_partners.delivery_partners(id);
-- ALTER TABLE orders.orders ADD COLUMN customer_name VARCHAR(100);
-- ALTER TABLE orders.orders ADD COLUMN customer_address_id UUID REFERENCES user_profile.addresses(id);

-- ALTER TABLE orders.orders ADD COLUMN pack_by_time TIMESTAMP;
-- ALTER TABLE orders.orders ADD COLUMN paid_by_time TIMESTAMP;
-- ALTER TABLE orders.orders ADD COLUMN delivered_by_time TIMESTAMP;

-- Ensure order_status_logs has proper structure
-- ALTER TABLE orders.order_status_logs ADD COLUMN IF NOT EXISTS order_id UUID REFERENCES orders.orders(id);
-- ALTER TABLE orders.order_status_logs ADD COLUMN IF NOT EXISTS status VARCHAR(50);
-- ALTER TABLE orders.order_status_logs ADD COLUMN IF NOT EXISTS changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;



-- Schema Changes for Inventory Management

-- 1. Add vendor inventory table to track items per vendor
-- CREATE TABLE vendor_inventory (
--                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--                                   vendor_id UUID NOT NULL REFERENCES vendor_accounts.vendor_accounts(id),
--                                   item_id INTEGER NOT NULL REFERENCES vendor_items.items(id),
--                                   stock_quantity INTEGER NOT NULL DEFAULT 0,
--                                   is_available BOOLEAN NOT NULL DEFAULT true,
--                                   price_override NUMERIC(10,2), -- Allow vendor to override item price
--                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                                   UNIQUE(vendor_id, item_id) -- Prevent duplicate items per vendor
-- );
--
-- -- 2. Add inventory tracking/audit table
-- CREATE TABLE inventory_movements (
--                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--                                      vendor_inventory_id UUID NOT NULL REFERENCES vendor_inventory(id),
--                                      movement_type VARCHAR(20) NOT NULL, -- 'add', 'remove', 'sold', 'expired', 'adjustment'
--                                      quantity_change INTEGER NOT NULL, -- positive for add, negative for remove
--                                      previous_quantity INTEGER NOT NULL,
--                                      new_quantity INTEGER NOT NULL,
--                                      reason VARCHAR(255),
--                                      created_by UUID REFERENCES user_profile.users(id), -- who made the change
--                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
--
-- -- 3. Update items table to support better search and categorization
-- ALTER TABLE vendor_items.items ADD COLUMN IF NOT EXISTS search_keywords TEXT; -- for better search
-- ALTER TABLE vendor_items.items ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
-- ALTER TABLE vendor_items.items ADD COLUMN IF NOT EXISTS vendor_id UUID REFERENCES vendor_accounts.vendor_accounts(id); -- if items are vendor-specific
--
-- -- 4. Indexes for performance
-- CREATE INDEX idx_vendor_inventory_vendor_id ON vendor_inventory(vendor_id);
-- CREATE INDEX idx_vendor_inventory_item_id ON vendor_inventory(item_id);
-- CREATE INDEX idx_vendor_inventory_available ON vendor_inventory(vendor_id, is_available);
-- CREATE INDEX idx_inventory_movements_vendor_inventory ON inventory_movements(vendor_inventory_id);
-- CREATE INDEX idx_items_category_active ON vendor_items.items(category_id, is_active);
-- CREATE INDEX idx_items_search ON vendor_items.items USING gin(to_tsvector('english', name || ' ' || COALESCE(description, '') || ' ' || COALESCE(search_keywords, '')));
--
-- -- 5. Create view for vendor inventory summary
-- CREATE OR REPLACE VIEW vendor_inventory_summary AS
-- SELECT
--     va.id as vendor_id,
--     COUNT(vi.id) as total_items,
--     COUNT(CASE WHEN vi.is_available = true AND vi.stock_quantity > 0 THEN 1 END) as in_stock_items,
--     COUNT(CASE WHEN vi.stock_quantity = 0 THEN 1 END) as out_of_stock_items
-- FROM vendor_accounts.vendor_accounts va
--          LEFT JOIN vendor_inventory vi ON va.id = vi.vendor_id
-- GROUP BY va.id;