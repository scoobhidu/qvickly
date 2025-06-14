--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.2 (Debian 17.2-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: delivery_partners; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA delivery_partners;


ALTER SCHEMA delivery_partners OWNER TO postgres;

--
-- Name: orders; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA orders;


ALTER SCHEMA orders OWNER TO postgres;

--
-- Name: qvickly_grocery_products; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA qvickly_grocery_products;


ALTER SCHEMA qvickly_grocery_products OWNER TO postgres;

--
-- Name: user_profile; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA user_profile;


ALTER SCHEMA user_profile OWNER TO postgres;

--
-- Name: user_sso; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA user_sso;


ALTER SCHEMA user_sso OWNER TO postgres;

--
-- Name: vendor_accounts; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA vendor_accounts;


ALTER SCHEMA vendor_accounts OWNER TO postgres;

--
-- Name: vendor_constants; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA vendor_constants;


ALTER SCHEMA vendor_constants OWNER TO postgres;

--
-- Name: vendor_items; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA vendor_items;


ALTER SCHEMA vendor_items OWNER TO postgres;

--
-- Name: account_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.account_type AS ENUM (
    'store',
    'restaurant'
);


ALTER TYPE public.account_type OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: delivery_partners; Type: TABLE; Schema: delivery_partners; Owner: postgres
--

CREATE TABLE delivery_partners.delivery_partners (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    phone_number character varying(15) NOT NULL,
    pin character varying(10),
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE delivery_partners.delivery_partners OWNER TO postgres;

--
-- Name: order_items; Type: TABLE; Schema: orders; Owner: postgres
--

CREATE TABLE orders.order_items (
    id integer NOT NULL,
    order_id integer NOT NULL,
    item_id integer NOT NULL,
    quantity integer NOT NULL,
    instructions text,
    is_checked boolean DEFAULT false,
    CONSTRAINT order_items_quantity_check CHECK ((quantity > 0))
);


ALTER TABLE orders.order_items OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: orders; Owner: postgres
--

CREATE SEQUENCE orders.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE orders.order_items_id_seq OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: orders; Owner: postgres
--

ALTER SEQUENCE orders.order_items_id_seq OWNED BY orders.order_items.id;


--
-- Name: order_status_logs; Type: TABLE; Schema: orders; Owner: postgres
--

CREATE TABLE orders.order_status_logs (
    id integer NOT NULL,
    order_id integer NOT NULL,
    status text NOT NULL,
    changed_at timestamp without time zone DEFAULT now(),
    CONSTRAINT order_status_logs_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'accepted'::text, 'packed'::text, 'ready'::text, 'completed'::text, 'rejected'::text])))
);


ALTER TABLE orders.order_status_logs OWNER TO postgres;

--
-- Name: order_status_logs_id_seq; Type: SEQUENCE; Schema: orders; Owner: postgres
--

CREATE SEQUENCE orders.order_status_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE orders.order_status_logs_id_seq OWNER TO postgres;

--
-- Name: order_status_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: orders; Owner: postgres
--

ALTER SEQUENCE orders.order_status_logs_id_seq OWNED BY orders.order_status_logs.id;


--
-- Name: orders; Type: TABLE; Schema: orders; Owner: postgres
--

CREATE TABLE orders.orders (
    id integer NOT NULL,
    account_id uuid NOT NULL,
    user_id uuid NOT NULL,
    order_time timestamp without time zone NOT NULL,
    total_amount numeric(10,2) NOT NULL,
    status text DEFAULT 'pending'::text,
    location text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    customer_id uuid,
    delivery_partner_id uuid,
    customer_name character varying(100),
    customer_address_id uuid,
    pack_by_time timestamp without time zone,
    paid_by_time timestamp without time zone,
    delivered_by_time timestamp without time zone,
    CONSTRAINT orders_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'accepted'::text, 'packed'::text, 'ready'::text, 'completed'::text, 'cancelled'::text, 'rejected'::text])))
);


ALTER TABLE orders.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: orders; Owner: postgres
--

CREATE SEQUENCE orders.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE orders.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: orders; Owner: postgres
--

ALTER SEQUENCE orders.orders_id_seq OWNED BY orders.orders.id;


--
-- Name: inventory_movements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventory_movements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vendor_inventory_id uuid NOT NULL,
    movement_type character varying(20) NOT NULL,
    quantity_change integer NOT NULL,
    previous_quantity integer NOT NULL,
    new_quantity integer NOT NULL,
    reason character varying(255),
    created_by uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.inventory_movements OWNER TO postgres;

--
-- Name: vendor_inventory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.vendor_inventory (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vendor_id uuid NOT NULL,
    item_id integer NOT NULL,
    stock_quantity integer DEFAULT 0 NOT NULL,
    is_available boolean DEFAULT true NOT NULL,
    price_override numeric(10,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.vendor_inventory OWNER TO postgres;

--
-- Name: vendor_accounts; Type: TABLE; Schema: vendor_accounts; Owner: postgres
--

CREATE TABLE vendor_accounts.vendor_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    phone_number character varying(15) NOT NULL,
    account_type public.account_type NOT NULL,
    business_name character varying(100) NOT NULL,
    owner_name character varying(100),
    email character varying(100),
    address text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    gstin_number character varying(35),
    opening_time time without time zone,
    closing_time time without time zone,
    image_url text NOT NULL,
    live_status boolean NOT NULL
);


ALTER TABLE vendor_accounts.vendor_accounts OWNER TO postgres;

--
-- Name: vendor_inventory_summary; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.vendor_inventory_summary AS
 SELECT va.id AS vendor_id,
    count(vi.id) AS total_items,
    count(
        CASE
            WHEN ((vi.is_available = true) AND (vi.stock_quantity > 0)) THEN 1
            ELSE NULL::integer
        END) AS in_stock_items,
    count(
        CASE
            WHEN (vi.stock_quantity = 0) THEN 1
            ELSE NULL::integer
        END) AS out_of_stock_items
   FROM (vendor_accounts.vendor_accounts va
     LEFT JOIN public.vendor_inventory vi ON ((va.id = vi.vendor_id)))
  GROUP BY va.id;


ALTER VIEW public.vendor_inventory_summary OWNER TO postgres;

--
-- Name: dashboard_posters; Type: TABLE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE TABLE qvickly_grocery_products.dashboard_posters (
    id integer NOT NULL,
    image_url text NOT NULL,
    link_url text,
    active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE qvickly_grocery_products.dashboard_posters OWNER TO postgres;

--
-- Name: dashboard_posters_id_seq; Type: SEQUENCE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE SEQUENCE qvickly_grocery_products.dashboard_posters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE qvickly_grocery_products.dashboard_posters_id_seq OWNER TO postgres;

--
-- Name: dashboard_posters_id_seq; Type: SEQUENCE OWNED BY; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER SEQUENCE qvickly_grocery_products.dashboard_posters_id_seq OWNED BY qvickly_grocery_products.dashboard_posters.id;


--
-- Name: item_groups; Type: TABLE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE TABLE qvickly_grocery_products.item_groups (
    id integer NOT NULL,
    name text NOT NULL,
    slug text,
    display_order integer DEFAULT 0
);


ALTER TABLE qvickly_grocery_products.item_groups OWNER TO postgres;

--
-- Name: item_groups_id_seq; Type: SEQUENCE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE SEQUENCE qvickly_grocery_products.item_groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE qvickly_grocery_products.item_groups_id_seq OWNER TO postgres;

--
-- Name: item_groups_id_seq; Type: SEQUENCE OWNED BY; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER SEQUENCE qvickly_grocery_products.item_groups_id_seq OWNED BY qvickly_grocery_products.item_groups.id;


--
-- Name: items; Type: TABLE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE TABLE qvickly_grocery_products.items (
    id integer NOT NULL,
    name text NOT NULL,
    image_url text,
    group_id integer,
    price numeric(10,2),
    is_available boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    rating numeric(2,1),
    description text,
    sustainability_facts text,
    CONSTRAINT items_rating_check CHECK (((rating >= (0)::numeric) AND (rating <= (5)::numeric)))
);


ALTER TABLE qvickly_grocery_products.items OWNER TO postgres;

--
-- Name: items_id_seq; Type: SEQUENCE; Schema: qvickly_grocery_products; Owner: postgres
--

CREATE SEQUENCE qvickly_grocery_products.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE qvickly_grocery_products.items_id_seq OWNER TO postgres;

--
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER SEQUENCE qvickly_grocery_products.items_id_seq OWNED BY qvickly_grocery_products.items.id;


--
-- Name: addresses; Type: TABLE; Schema: user_profile; Owner: postgres
--

CREATE TABLE user_profile.addresses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    label character varying(50),
    address_line1 text NOT NULL,
    address_line2 text,
    city character varying(100),
    state character varying(100) NOT NULL,
    postal_code character varying(20) NOT NULL,
    country character varying(100) DEFAULT 'India'::character varying,
    latitude numeric(9,6) NOT NULL,
    longitude numeric(9,6) NOT NULL,
    is_default boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE user_profile.addresses OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: user_profile; Owner: postgres
--

CREATE TABLE user_profile.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    phone_number character varying(15) NOT NULL,
    email character varying(100) NOT NULL,
    full_name character varying(100),
    google_id character varying(255),
    profile_picture_url text,
    is_marketing_opted boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE user_profile.users OWNER TO postgres;

--
-- Name: user_sessions; Type: TABLE; Schema: user_sso; Owner: postgres
--

CREATE TABLE user_sso.user_sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    refresh_token text NOT NULL,
    refresh_token_expires_at timestamp without time zone NOT NULL,
    ip_address text,
    user_agent text,
    device_info jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_seen_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE user_sso.user_sessions OWNER TO postgres;

--
-- Name: restaurant_details; Type: TABLE; Schema: vendor_accounts; Owner: postgres
--

CREATE TABLE vendor_accounts.restaurant_details (
    account_id uuid NOT NULL,
    fssai_license_no character varying(50),
    cuisine_id uuid NOT NULL
);


ALTER TABLE vendor_accounts.restaurant_details OWNER TO postgres;

--
-- Name: cuisines; Type: TABLE; Schema: vendor_constants; Owner: postgres
--

CREATE TABLE vendor_constants.cuisines (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE vendor_constants.cuisines OWNER TO postgres;

--
-- Name: categories; Type: TABLE; Schema: vendor_items; Owner: postgres
--

CREATE TABLE vendor_items.categories (
    id integer NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE vendor_items.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: vendor_items; Owner: postgres
--

CREATE SEQUENCE vendor_items.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE vendor_items.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: vendor_items; Owner: postgres
--

ALTER SEQUENCE vendor_items.categories_id_seq OWNED BY vendor_items.categories.id;


--
-- Name: item_images; Type: TABLE; Schema: vendor_items; Owner: postgres
--

CREATE TABLE vendor_items.item_images (
    id integer NOT NULL,
    item_id integer NOT NULL,
    image_url text NOT NULL,
    "position" integer,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT item_images_position_check CHECK ((("position" >= 1) AND ("position" <= 4)))
);


ALTER TABLE vendor_items.item_images OWNER TO postgres;

--
-- Name: item_images_id_seq; Type: SEQUENCE; Schema: vendor_items; Owner: postgres
--

CREATE SEQUENCE vendor_items.item_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE vendor_items.item_images_id_seq OWNER TO postgres;

--
-- Name: item_images_id_seq; Type: SEQUENCE OWNED BY; Schema: vendor_items; Owner: postgres
--

ALTER SEQUENCE vendor_items.item_images_id_seq OWNED BY vendor_items.item_images.id;


--
-- Name: items; Type: TABLE; Schema: vendor_items; Owner: postgres
--

CREATE TABLE vendor_items.items (
    id integer NOT NULL,
    account_id uuid NOT NULL,
    category_id integer,
    name text NOT NULL,
    description text,
    price_retail numeric(10,2),
    price_wholesale numeric(10,2),
    is_available boolean DEFAULT true,
    stock integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    search_keywords text,
    is_active boolean DEFAULT true,
    vendor_id uuid
);


ALTER TABLE vendor_items.items OWNER TO postgres;

--
-- Name: items_id_seq; Type: SEQUENCE; Schema: vendor_items; Owner: postgres
--

CREATE SEQUENCE vendor_items.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE vendor_items.items_id_seq OWNER TO postgres;

--
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: vendor_items; Owner: postgres
--

ALTER SEQUENCE vendor_items.items_id_seq OWNED BY vendor_items.items.id;


--
-- Name: order_items id; Type: DEFAULT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_items ALTER COLUMN id SET DEFAULT nextval('orders.order_items_id_seq'::regclass);


--
-- Name: order_status_logs id; Type: DEFAULT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_status_logs ALTER COLUMN id SET DEFAULT nextval('orders.order_status_logs_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders ALTER COLUMN id SET DEFAULT nextval('orders.orders_id_seq'::regclass);


--
-- Name: dashboard_posters id; Type: DEFAULT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.dashboard_posters ALTER COLUMN id SET DEFAULT nextval('qvickly_grocery_products.dashboard_posters_id_seq'::regclass);


--
-- Name: item_groups id; Type: DEFAULT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.item_groups ALTER COLUMN id SET DEFAULT nextval('qvickly_grocery_products.item_groups_id_seq'::regclass);


--
-- Name: items id; Type: DEFAULT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.items ALTER COLUMN id SET DEFAULT nextval('qvickly_grocery_products.items_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.categories ALTER COLUMN id SET DEFAULT nextval('vendor_items.categories_id_seq'::regclass);


--
-- Name: item_images id; Type: DEFAULT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.item_images ALTER COLUMN id SET DEFAULT nextval('vendor_items.item_images_id_seq'::regclass);


--
-- Name: items id; Type: DEFAULT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.items ALTER COLUMN id SET DEFAULT nextval('vendor_items.items_id_seq'::regclass);


--
-- Data for Name: delivery_partners; Type: TABLE DATA; Schema: delivery_partners; Owner: postgres
--

COPY delivery_partners.delivery_partners (id, name, phone_number, pin, is_active, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: orders; Owner: postgres
--

COPY orders.order_items (id, order_id, item_id, quantity, instructions, is_checked) FROM stdin;
\.


--
-- Data for Name: order_status_logs; Type: TABLE DATA; Schema: orders; Owner: postgres
--

COPY orders.order_status_logs (id, order_id, status, changed_at) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: orders; Owner: postgres
--

COPY orders.orders (id, account_id, user_id, order_time, total_amount, status, location, created_at, updated_at, customer_id, delivery_partner_id, customer_name, customer_address_id, pack_by_time, paid_by_time, delivered_by_time) FROM stdin;
1	69fc8ba4-d11f-4618-9042-c1523d381013	aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa	2025-05-22 05:30:03.348509	149.99	pending	Delhi	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
3	69fc8ba4-d11f-4618-9042-c1523d381013	cccccccc-cccc-cccc-cccc-cccccccccccc	2025-05-24 02:30:03.348509	89.00	packed	Jharkhand	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
5	69fc8ba4-d11f-4618-9042-c1523d381013	eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee	2025-05-24 04:30:03.348509	300.00	completed	Austin, TX	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
6	69fc8ba4-d11f-4618-9042-c1523d381013	ffffffff-ffff-ffff-ffff-ffffffffffff	2025-05-26 05:30:03.348	200.00	cancelled	Boston, MA	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
4	69fc8ba4-d11f-4618-9042-c1523d381013	dddddddd-dddd-dddd-dddd-dddddddddddd	2025-05-26 23:30:03.348	123.45	ready	Manali	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
2	69fc8ba4-d11f-4618-9042-c1523d381013	bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb	2025-05-26 05:30:03.348	259.75	accepted	Ambikapur	2025-05-24 05:30:03.348509	2025-05-24 05:30:03.348509	\N	\N	\N	\N	\N	\N	\N
\.


--
-- Data for Name: inventory_movements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventory_movements (id, vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_by, created_at) FROM stdin;
\.


--
-- Data for Name: vendor_inventory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.vendor_inventory (id, vendor_id, item_id, stock_quantity, is_available, price_override, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: dashboard_posters; Type: TABLE DATA; Schema: qvickly_grocery_products; Owner: postgres
--

COPY qvickly_grocery_products.dashboard_posters (id, image_url, link_url, active, created_at) FROM stdin;
\.


--
-- Data for Name: item_groups; Type: TABLE DATA; Schema: qvickly_grocery_products; Owner: postgres
--

COPY qvickly_grocery_products.item_groups (id, name, slug, display_order) FROM stdin;
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: qvickly_grocery_products; Owner: postgres
--

COPY qvickly_grocery_products.items (id, name, image_url, group_id, price, is_available, created_at, rating, description, sustainability_facts) FROM stdin;
\.


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: user_profile; Owner: postgres
--

COPY user_profile.addresses (id, user_id, label, address_line1, address_line2, city, state, postal_code, country, latitude, longitude, is_default, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: user_profile; Owner: postgres
--

COPY user_profile.users (id, phone_number, email, full_name, google_id, profile_picture_url, is_marketing_opted, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: user_sessions; Type: TABLE DATA; Schema: user_sso; Owner: postgres
--

COPY user_sso.user_sessions (id, user_id, refresh_token, refresh_token_expires_at, ip_address, user_agent, device_info, created_at, last_seen_at) FROM stdin;
\.


--
-- Data for Name: restaurant_details; Type: TABLE DATA; Schema: vendor_accounts; Owner: postgres
--

COPY vendor_accounts.restaurant_details (account_id, fssai_license_no, cuisine_id) FROM stdin;
\.


--
-- Data for Name: vendor_accounts; Type: TABLE DATA; Schema: vendor_accounts; Owner: postgres
--

COPY vendor_accounts.vendor_accounts (id, phone_number, account_type, business_name, owner_name, email, address, created_at, latitude, longitude, gstin_number, opening_time, closing_time, image_url, live_status) FROM stdin;
69fc8ba4-d11f-4618-9042-c1523d381013	8010201921	restaurant	abc	rajat	rajatnd9@gmail.com	awdawdawdawd	2025-05-26 15:47:00.822	12.3445245	23.45512	123AwWXAAWXAWVDF123123	15:47:00.822166	15:47:00.822166	s3_url	t
\.


--
-- Data for Name: cuisines; Type: TABLE DATA; Schema: vendor_constants; Owner: postgres
--

COPY vendor_constants.cuisines (id, name) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: vendor_items; Owner: postgres
--

COPY vendor_items.categories (id, name, created_at) FROM stdin;
\.


--
-- Data for Name: item_images; Type: TABLE DATA; Schema: vendor_items; Owner: postgres
--

COPY vendor_items.item_images (id, item_id, image_url, "position", created_at) FROM stdin;
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: vendor_items; Owner: postgres
--

COPY vendor_items.items (id, account_id, category_id, name, description, price_retail, price_wholesale, is_available, stock, created_at, updated_at, search_keywords, is_active, vendor_id) FROM stdin;
\.


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: orders; Owner: postgres
--

SELECT pg_catalog.setval('orders.order_items_id_seq', 1, false);


--
-- Name: order_status_logs_id_seq; Type: SEQUENCE SET; Schema: orders; Owner: postgres
--

SELECT pg_catalog.setval('orders.order_status_logs_id_seq', 1, false);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: orders; Owner: postgres
--

SELECT pg_catalog.setval('orders.orders_id_seq', 6, true);


--
-- Name: dashboard_posters_id_seq; Type: SEQUENCE SET; Schema: qvickly_grocery_products; Owner: postgres
--

SELECT pg_catalog.setval('qvickly_grocery_products.dashboard_posters_id_seq', 1, false);


--
-- Name: item_groups_id_seq; Type: SEQUENCE SET; Schema: qvickly_grocery_products; Owner: postgres
--

SELECT pg_catalog.setval('qvickly_grocery_products.item_groups_id_seq', 1, false);


--
-- Name: items_id_seq; Type: SEQUENCE SET; Schema: qvickly_grocery_products; Owner: postgres
--

SELECT pg_catalog.setval('qvickly_grocery_products.items_id_seq', 1, false);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: vendor_items; Owner: postgres
--

SELECT pg_catalog.setval('vendor_items.categories_id_seq', 1, false);


--
-- Name: item_images_id_seq; Type: SEQUENCE SET; Schema: vendor_items; Owner: postgres
--

SELECT pg_catalog.setval('vendor_items.item_images_id_seq', 1, false);


--
-- Name: items_id_seq; Type: SEQUENCE SET; Schema: vendor_items; Owner: postgres
--

SELECT pg_catalog.setval('vendor_items.items_id_seq', 1, false);


--
-- Name: delivery_partners delivery_partners_pkey; Type: CONSTRAINT; Schema: delivery_partners; Owner: postgres
--

ALTER TABLE ONLY delivery_partners.delivery_partners
    ADD CONSTRAINT delivery_partners_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: order_status_logs order_status_logs_pkey; Type: CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_status_logs
    ADD CONSTRAINT order_status_logs_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: inventory_movements inventory_movements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_pkey PRIMARY KEY (id);


--
-- Name: vendor_inventory vendor_inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vendor_inventory
    ADD CONSTRAINT vendor_inventory_pkey PRIMARY KEY (id);


--
-- Name: vendor_inventory vendor_inventory_vendor_id_item_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vendor_inventory
    ADD CONSTRAINT vendor_inventory_vendor_id_item_id_key UNIQUE (vendor_id, item_id);


--
-- Name: dashboard_posters dashboard_posters_pkey; Type: CONSTRAINT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.dashboard_posters
    ADD CONSTRAINT dashboard_posters_pkey PRIMARY KEY (id);


--
-- Name: item_groups item_groups_pkey; Type: CONSTRAINT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.item_groups
    ADD CONSTRAINT item_groups_pkey PRIMARY KEY (id);


--
-- Name: item_groups item_groups_slug_key; Type: CONSTRAINT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.item_groups
    ADD CONSTRAINT item_groups_slug_key UNIQUE (slug);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_google_id_key; Type: CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.users
    ADD CONSTRAINT users_google_id_key UNIQUE (google_id);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: user_sso; Owner: postgres
--

ALTER TABLE ONLY user_sso.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_refresh_token_key; Type: CONSTRAINT; Schema: user_sso; Owner: postgres
--

ALTER TABLE ONLY user_sso.user_sessions
    ADD CONSTRAINT user_sessions_refresh_token_key UNIQUE (refresh_token);


--
-- Name: restaurant_details acc_cuis_key; Type: CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.restaurant_details
    ADD CONSTRAINT acc_cuis_key UNIQUE (account_id, cuisine_id);


--
-- Name: restaurant_details restaurant_details_pkey; Type: CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.restaurant_details
    ADD CONSTRAINT restaurant_details_pkey PRIMARY KEY (account_id);


--
-- Name: vendor_accounts vendor_accounts_phone_number_key; Type: CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.vendor_accounts
    ADD CONSTRAINT vendor_accounts_phone_number_key UNIQUE (phone_number);


--
-- Name: vendor_accounts vendor_accounts_pkey; Type: CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.vendor_accounts
    ADD CONSTRAINT vendor_accounts_pkey PRIMARY KEY (id);


--
-- Name: cuisines cuisines_name_key; Type: CONSTRAINT; Schema: vendor_constants; Owner: postgres
--

ALTER TABLE ONLY vendor_constants.cuisines
    ADD CONSTRAINT cuisines_name_key UNIQUE (name);


--
-- Name: cuisines cuisines_pkey; Type: CONSTRAINT; Schema: vendor_constants; Owner: postgres
--

ALTER TABLE ONLY vendor_constants.cuisines
    ADD CONSTRAINT cuisines_pkey PRIMARY KEY (id);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: item_images item_images_pkey; Type: CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.item_images
    ADD CONSTRAINT item_images_pkey PRIMARY KEY (id);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: idx_inventory_movements_vendor_inventory; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inventory_movements_vendor_inventory ON public.inventory_movements USING btree (vendor_inventory_id);


--
-- Name: idx_vendor_inventory_available; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_vendor_inventory_available ON public.vendor_inventory USING btree (vendor_id, is_available);


--
-- Name: idx_vendor_inventory_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_vendor_inventory_item_id ON public.vendor_inventory USING btree (item_id);


--
-- Name: idx_vendor_inventory_vendor_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_vendor_inventory_vendor_id ON public.vendor_inventory USING btree (vendor_id);


--
-- Name: idx_item_images_item_id; Type: INDEX; Schema: vendor_items; Owner: postgres
--

CREATE INDEX idx_item_images_item_id ON vendor_items.item_images USING btree (item_id);


--
-- Name: idx_items_account_id; Type: INDEX; Schema: vendor_items; Owner: postgres
--

CREATE INDEX idx_items_account_id ON vendor_items.items USING btree (account_id);


--
-- Name: idx_items_category_active; Type: INDEX; Schema: vendor_items; Owner: postgres
--

CREATE INDEX idx_items_category_active ON vendor_items.items USING btree (category_id, is_active);


--
-- Name: idx_items_category_id; Type: INDEX; Schema: vendor_items; Owner: postgres
--

CREATE INDEX idx_items_category_id ON vendor_items.items USING btree (category_id);


--
-- Name: idx_items_search; Type: INDEX; Schema: vendor_items; Owner: postgres
--

CREATE INDEX idx_items_search ON vendor_items.items USING gin (to_tsvector('english'::regconfig, ((((name || ' '::text) || COALESCE(description, ''::text)) || ' '::text) || COALESCE(search_keywords, ''::text))));


--
-- Name: order_items order_items_item_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_items
    ADD CONSTRAINT order_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES vendor_items.items(id);


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders.orders(id) ON DELETE CASCADE;


--
-- Name: order_status_logs order_status_logs_order_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.order_status_logs
    ADD CONSTRAINT order_status_logs_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders.orders(id) ON DELETE CASCADE;


--
-- Name: orders orders_account_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders
    ADD CONSTRAINT orders_account_id_fkey FOREIGN KEY (account_id) REFERENCES vendor_accounts.vendor_accounts(id) ON DELETE CASCADE;


--
-- Name: orders orders_customer_address_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders
    ADD CONSTRAINT orders_customer_address_id_fkey FOREIGN KEY (customer_address_id) REFERENCES user_profile.addresses(id);


--
-- Name: orders orders_delivery_partner_id_fkey; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders
    ADD CONSTRAINT orders_delivery_partner_id_fkey FOREIGN KEY (delivery_partner_id) REFERENCES delivery_partners.delivery_partners(id);


--
-- Name: orders orders_orders__fk; Type: FK CONSTRAINT; Schema: orders; Owner: postgres
--

ALTER TABLE ONLY orders.orders
    ADD CONSTRAINT orders_orders__fk FOREIGN KEY (customer_id) REFERENCES user_profile.users(id);


--
-- Name: inventory_movements inventory_movements_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_created_by_fkey FOREIGN KEY (created_by) REFERENCES user_profile.users(id);


--
-- Name: inventory_movements inventory_movements_vendor_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_vendor_inventory_id_fkey FOREIGN KEY (vendor_inventory_id) REFERENCES public.vendor_inventory(id);


--
-- Name: vendor_inventory vendor_inventory_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vendor_inventory
    ADD CONSTRAINT vendor_inventory_item_id_fkey FOREIGN KEY (item_id) REFERENCES vendor_items.items(id);


--
-- Name: vendor_inventory vendor_inventory_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vendor_inventory
    ADD CONSTRAINT vendor_inventory_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES vendor_accounts.vendor_accounts(id);


--
-- Name: items items_group_id_fkey; Type: FK CONSTRAINT; Schema: qvickly_grocery_products; Owner: postgres
--

ALTER TABLE ONLY qvickly_grocery_products.items
    ADD CONSTRAINT items_group_id_fkey FOREIGN KEY (group_id) REFERENCES qvickly_grocery_products.item_groups(id) ON DELETE CASCADE;


--
-- Name: addresses addresses_user_id_fkey; Type: FK CONSTRAINT; Schema: user_profile; Owner: postgres
--

ALTER TABLE ONLY user_profile.addresses
    ADD CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profile.users(id) ON DELETE CASCADE;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: user_sso; Owner: postgres
--

ALTER TABLE ONLY user_sso.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profile.users(id) ON DELETE CASCADE;


--
-- Name: restaurant_details restaurant_details_account_id_fkey; Type: FK CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.restaurant_details
    ADD CONSTRAINT restaurant_details_account_id_fkey FOREIGN KEY (account_id) REFERENCES vendor_accounts.vendor_accounts(id) ON DELETE CASCADE;


--
-- Name: restaurant_details restaurant_details_cuisine_id_fkey; Type: FK CONSTRAINT; Schema: vendor_accounts; Owner: postgres
--

ALTER TABLE ONLY vendor_accounts.restaurant_details
    ADD CONSTRAINT restaurant_details_cuisine_id_fkey FOREIGN KEY (cuisine_id) REFERENCES vendor_constants.cuisines(id) ON DELETE CASCADE;


--
-- Name: item_images item_images_item_id_fkey; Type: FK CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.item_images
    ADD CONSTRAINT item_images_item_id_fkey FOREIGN KEY (item_id) REFERENCES vendor_items.items(id) ON DELETE CASCADE;


--
-- Name: items items_account_id_fkey; Type: FK CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.items
    ADD CONSTRAINT items_account_id_fkey FOREIGN KEY (account_id) REFERENCES vendor_accounts.vendor_accounts(id) ON DELETE CASCADE;


--
-- Name: items items_category_id_fkey; Type: FK CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.items
    ADD CONSTRAINT items_category_id_fkey FOREIGN KEY (category_id) REFERENCES vendor_items.categories(id);


--
-- Name: items items_vendor_id_fkey; Type: FK CONSTRAINT; Schema: vendor_items; Owner: postgres
--

ALTER TABLE ONLY vendor_items.items
    ADD CONSTRAINT items_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES vendor_accounts.vendor_accounts(id);


--
-- PostgreSQL database dump complete
--


ALTER table vendor_accounts.vendor_accounts
    add password varchar(20) default substr(md5(random()::text), 0, random(8, 20)) not null;


