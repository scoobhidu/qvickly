CREATE DATABASE quickkart;



CREATE SCHEMA IF NOT EXISTS master;
CREATE SCHEMA IF NOT EXISTS profile;
CREATE SCHEMA IF NOT EXISTS customer;
CREATE SCHEMA IF NOT EXISTS delivery;
CREATE SCHEMA IF NOT EXISTS vendor;

-- ============================================
-- MASTER SCHEMA TABLES (Independent tables)
-- ============================================

-- Create sequence first (used by grocery_items)
create sequence items_item_id_seq;
alter sequence items_item_id_seq owner to postgres;

-- Cuisines table (referenced by vendors)
create table if not exists master.cuisines
(
    cuisine_id uuid default gen_random_uuid() not null
        primary key,
    title      varchar                        not null,
    image_url  text default ''::text
);

alter table master.cuisines
    owner to postgres;

-- Grocery categories (parent of subcategories)
create table if not exists master.grocery_categories
(
    grocery_category_id bigserial
        primary key,
    title               varchar(50) not null,
    created_at          timestamp default now()
);

alter table master.grocery_categories
    owner to postgres;

-- Grocery subcategories (child of categories, parent of items)
create table if not exists master.grocery_subcategories
(
    grocery_subcategory_id bigserial
        primary key,
    grocery_category_id    bigint
        references master.grocery_categories,
    title                  varchar(50) not null,
    created_at             timestamp default now()
);

alter table master.grocery_subcategories
    owner to postgres;

-- Grocery items (child of categories and subcategories)
create table if not exists master.grocery_items
(
    item_id         bigint    default nextval('master.items_item_id_seq'::regclass) not null
        constraint items_pkey
            primary key,
    category_id     bigint
        constraint items_category_id_fkey
            references master.grocery_categories,
    subcategory_id  bigint
        constraint items_subcategory_id_fkey
            references master.grocery_subcategories,
    title           text                                                            not null,
    description     text      default ''::text,
    price_retail    numeric(10, 2)                                                  not null,
    price_wholesale numeric(10, 2)                                                  not null,
    search_keywords text[]    default '{}'::text[],
    created_at      timestamp default now(),
    image_url_1     text                                                            not null,
    image_url_2     text                                                            not null,
    image_url_3     text                                                            not null,
    image_url_4     text                                                            not null
);

alter table master.grocery_items
    owner to postgres;

alter sequence items_item_id_seq owned by master.grocery_items.item_id;

-- ============================================
-- PROFILE SCHEMA TABLES
-- ============================================

-- Customer table (independent, referenced by addresses and orders)
create table if not exists profile.customer
(
    id                 uuid      default gen_random_uuid()     not null
        primary key,
    phone              varchar(10)                             not null,
    email              varchar(50)                             not null,
    full_name          varchar,
    google_id          varchar,
    profile_pic        varchar,
    is_marketing_opted boolean   default true,
    created_at         timestamp default now(),
    updated_at         timestamp default now(),
    password           varchar   default ''::character varying not null
);

alter table profile.customer
    owner to postgres;

-- Customer addresses (child of customer)
create table if not exists profile.customer_addresses
(
    address_id    bigserial
        primary key,
    customer_id   uuid
        references profile.customer,
    title         varchar(50)    not null,
    address_line1 text           not null,
    address_line2 text,
    city          varchar        not null,
    state         varchar        not null,
    postal_code   varchar(6)     not null,
    country       varchar   default 'India'::character varying,
    latitude      numeric(10, 7) not null,
    longitude     numeric(11, 7) not null,
    is_default    boolean   default false,
    created_at    timestamp default now(),
    updated_at    timestamp default now()
);

alter table profile.customer_addresses
    owner to postgres;

-- Delivery boy table (independent, referenced by order_tracker)
create table if not exists profile.delivery_boy
(
    id                      uuid      default gen_random_uuid()     not null
        primary key,
    password                varchar                                 not null,
    full_name               varchar                                 not null,
    phone                   varchar(10)                             not null,
    aadhar                  bigint                                  not null,
    pin                     bigint    default random(100000, 999999),
    is_active               boolean   default true,
    image_url               varchar   default ''::character varying not null,
    latitude                numeric(10, 7)                          not null,
    longitude               numeric(11, 7)                          not null,
    is_marketing_opted      boolean   default true,
    created_at              timestamp default now(),
    updated_at              timestamp default now(),
    rating                  integer   default 5,
    upi_id                  text      default ''::text              not null,
    emergency_contact_phone varchar(10)                             not null,
    emergency_contact_name  text                                    not null,
    date_of_birth           timestamp
);

alter table profile.delivery_boy
    owner to postgres;

-- Vendors table (child of cuisines, parent of inventory and order_pickup_assignments)
create table if not exists profile.vendors
(
    vendor_id     uuid      default gen_random_uuid()     not null
        primary key,
    password      varchar                                 not null,
    account_type  text
        constraint vendors_account_type_check
            check (lower(account_type) = ANY (ARRAY ['grocery'::text, 'restaurant'::text])),
    business_name varchar                                 not null,
    owner_name    varchar                                 not null,
    phone         varchar(10)                             not null,
    aadhar        bigint                                  not null,
    is_active     boolean   default true,
    image_url     varchar   default ''::character varying not null,
    gstin_number  varchar(20)                             not null,
    opening_time  time                                    not null,
    closing_time  time                                    not null,
    address       text                                    not null,
    latitude      numeric(10, 7)                          not null,
    longitude     numeric(11, 7)                          not null,
    created_at    timestamp default now(),
    updated_at    timestamp default now(),
    cuisine_id    uuid
        references master.cuisines
);

alter table profile.vendors
    owner to postgres;

-- ============================================
-- CUSTOMER SCHEMA TABLES (Orders and related)
-- ============================================

-- Orders table (child of customer, parent of order_items and order_pickup_assignments)
create table if not exists customer.orders
(
    order_id        uuid      default gen_random_uuid() not null
        primary key,
    customer_id     uuid
        references profile.customer,
    order_time      timestamp default now(),
    pack_by_time    timestamp default (now() + '00:08:00'::interval),
    pick_up_time    timestamp default (now() + '00:10:00'::interval),
    deliver_by_time timestamp default (now() + '00:30:00'::interval),
    paid_time       timestamp,
    delivery_time   timestamp,
    updated_at      timestamp default now(),
    instructions    varchar   default ''::character varying,
    amount          double precision,
    status          text                                not null
        constraint orders_status_check
            check (lower(status) = ANY
                   (ARRAY ['pending'::text, 'accepted'::text, 'packed'::text, 'picked'::text, 'delivered'::text, 'cancelled'::text, 'rejected'::text]))
);

alter table customer.orders
    owner to postgres;

-- Order items table (child of orders and grocery_items) - FIRST VERSION
create table if not exists customer.order_items
(
    order_id uuid    not null
        references customer.orders,
    item_id  bigint  not null
        references master.grocery_items,
    qty      integer not null,
    primary key (order_id, item_id)
);

alter table customer.order_items
    owner to postgres;

-- Order tracker table (child of orders and delivery_boy)
create table if not exists customer.order_tracker
(
    delivery_id     uuid default gen_random_uuid() not null
        constraint order_tracker_pk
            primary key,
    order_id        uuid
        references customer.orders,
    delivery_boy_id uuid
        references profile.delivery_boy
);

alter table customer.order_tracker
    owner to postgres;

-- ============================================
-- VENDOR SCHEMA TABLES
-- ============================================

-- Inventory table (child of grocery_items and vendors)
create table if not exists vendor.inventory
(
    item_id                  bigserial
        references master.grocery_items,
    vendor_id                uuid
        references profile.vendors,
    created_at               timestamp default now(),
    updated_at               timestamp default now(),
    qty                      integer   default 0 not null,
    wholesale_price_override double precision
);

alter table vendor.inventory
    owner to postgres;

-- Order pickup assignments (child of vendors and orders)
create table if not exists vendor.order_pickup_assignments
(
    vendor_assignment_id uuid    default gen_random_uuid() not null
        primary key,
    vendor_id            uuid
        references profile.vendors,
    order_id             uuid
        references customer.orders,
    picked_up            boolean default false
);

alter table vendor.order_pickup_assignments
    owner to postgres;

-- Order items table (child of order_pickup_assignments and grocery_items) - SECOND VERSION
create table if not exists vendor.order_items
(
    vendor_assignment_id uuid    not null
        references vendor.order_pickup_assignments,
    item_id              bigint  not null
        references master.grocery_items,
    qty                  integer not null,
    primary key (vendor_assignment_id, item_id)
);

alter table vendor.order_items
    owner to postgres;

-- Vendor pickup tracker (child of order_tracker and order_pickup_assignments)
create table if not exists vendor.vendor_pickup_tracker
(
    delivery_id          uuid not null
        references customer.order_tracker,
    vendor_assignment_id uuid not null
        references vendor.order_pickup_assignments,
    delivery_boy_id      uuid
        references profile.delivery_boy,
    primary key (delivery_id, vendor_assignment_id)
);

alter table vendor.vendor_pickup_tracker
    owner to postgres;

create table if not exists delivery.order_tracker
(
    delivery_id     uuid default gen_random_uuid() not null
        constraint order_tracker_pk
            primary key,
    order_id        uuid
        references customer.orders,
    delivery_boy_id uuid
        references profile.delivery_boy
);

alter table delivery.order_tracker
    owner to postgres;

create table if not exists delivery.vendor_pickup_tracker
(
    delivery_id          uuid not null
        references delivery.order_tracker,
    vendor_assignment_id uuid not null
        references vendor.order_pickup_assignments,
    delivery_boy_id      uuid
        references profile.delivery_boy,
    primary key (delivery_id, vendor_assignment_id)
);

alter table delivery.vendor_pickup_tracker
    owner to postgres;


INSERT INTO master.grocery_categories (title, created_at)
VALUES ('Tea, Coffee, Juice', NOW());

-- 2. Insert into grocery_subcategories (depends on grocery_categories)
INSERT INTO master.grocery_subcategories (grocery_category_id, title, created_at)
VALUES (1, 'Tea', NOW());

-- 3. Insert into cuisines (independent table)
INSERT INTO master.cuisines (title, image_url)
VALUES ('Italian', 'https://example.com/italian-cuisine.jpg');

-- 4. Insert into customer (independent table)
INSERT INTO profile.customer (phone, email, full_name, google_id, profile_pic, is_marketing_opted, created_at, updated_at, password)
VALUES ('9876543210', 'john.doe@email.com', 'John Doe', 'google_123456', 'https://example.com/profile.jpg', true, NOW(), NOW(), 'hashed_password_123');

-- 5. Insert into customer_addresses (depends on customer)
INSERT INTO profile.customer_addresses (customer_id, title, address_line1, address_line2, city, state, postal_code, country, latitude, longitude, is_default, created_at, updated_at)
VALUES ('afdfbb3f-bbea-43f8-86fb-9bd4d86e6c46'::uuid, 'Home', '123 Main Street', 'Apt 4B', 'New York', 'NY', '10001', 'USA', 9.7128, -9.0060, true, NOW(), NOW());

-- 6. Insert into vendors (independent table)
INSERT INTO profile.vendors (password, account_type, business_name, owner_name, phone, is_active, image_url, gstin_number, aadhar, opening_time, closing_time, address, latitude, longitude, created_at, updated_at, cuisine_id)
VALUES ('vendor_password_123', 'restaurant', 'Marios Pizza', 'Mario Rossi', '9876543211', true, 'https://example.com/vendor.jpg', 'GST123456789', '567890', '09:00:00', '22:00:00', '456 Restaurant Ave, New York, NY', 8.7589, -8.9851, NOW(), NOW(), 'bc9fdd7d-72ff-4a64-82d0-0eb06bebc665'::uuid);

-- 7. Insert into delivery_boy (independent table)
INSERT INTO profile.delivery_boy (password, full_name, phone, aadhar, pin, is_active, image_url, latitude, longitude, is_marketing_opted, created_at, updated_at, rating, upi_id, emergency_contact_phone, emergency_contact_name, date_of_birth)
VALUES ('delivery_password_123', 'Raj Kumar', '9876543212', '123456789012', '654321', true, 'https://example.com/delivery.jpg', 9.7505, -9.9934, false, NOW(), NOW(), 5, 'raj@paytm', '9876543213', 'Emergency Contact', '1990-05-05');

-- 8. Insert into grocery_items (depends on grocery_subcategories)
INSERT INTO master.grocery_items (category_id, subcategory_id, title, description, price_retail, price_wholesale, search_keywords, created_at, image_url_1, image_url_2, image_url_3, image_url_4)
VALUES (1, 1, 'Fresh Apples', 'Red delicious apples from local farms', 5.99, 4.50, array['apples', 'fruit', 'fresh', 'red'], NOW(), 'https://example.com/apple1.jpg', 'https://example.com/apple2.jpg', 'https://example.com/apple3.jpg', 'https://example.com/apple4.jpg');

-- 9. Insert into inventory (depends on vendors and grocery_items)
INSERT INTO vendor.inventory (item_id, vendor_id, created_at, updated_at, qty, wholesale_price_override)
VALUES (1, 'b001cd89-d713-4c26-af01-87b07611792e'::uuid, NOW(), NOW(), 100, 4.25);

-- 10. Insert into orders (depends on customer)
INSERT INTO customer.orders (customer_id, order_time, pack_by_time, pick_up_time, deliver_by_time, paid_time, delivery_time, updated_at, instructions, amount, status)
VALUES ('afdfbb3f-bbea-43f8-86fb-9bd4d86e6c46'::uuid, NOW(), NOW() + INTERVAL '1 hour', NOW() + INTERVAL '2 hours', NOW() + INTERVAL '3 hours', NOW(), NOW() + INTERVAL '3 hours', NOW(), 'Please handle with care', 25.99, 'delivered');

-- 11. Insert into order_items (depends on orders and grocery_items)
INSERT INTO customer.order_items (qty, order_id, item_id)
VALUES (4, 'b2a4f55a-b5b2-47b4-b3a5-4a68583aed2b'::uuid, 1);

-- 12. Insert into order_pickup_assignments (depends on orders and vendors)
INSERT INTO vendor.order_pickup_assignments (vendor_id, order_id, picked_up)
VALUES ('b001cd89-d713-4c26-af01-87b07611792e'::uuid, 'b2a4f55a-b5b2-47b4-b3a5-4a68583aed2b', true);

-- 13. Insert into order_tracker (depends on orders and delivery_boy)
INSERT INTO delivery.order_tracker (order_id)
VALUES ('b2a4f55a-b5b2-47b4-b3a5-4a68583aed2b'::uuid);

-- 14. Insert into vendor_pickup_tracker (depends on delivery_boy and vendor_assignment)
INSERT INTO delivery.vendor_pickup_tracker (delivery_id, vendor_assignment_id)
VALUES ('0eb2156d-57c7-4723-89f5-c54295683891'::uuid, '333e7532-dc80-4725-9c63-ed6ac2b2ee23'::uuid);