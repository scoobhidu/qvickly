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


-- Clear existing data (if any)
TRUNCATE TABLE orders.order_items CASCADE;
TRUNCATE TABLE orders.order_status_logs CASCADE;
TRUNCATE TABLE orders.orders CASCADE;
TRUNCATE TABLE public.inventory_movements CASCADE;
TRUNCATE TABLE public.vendor_inventory CASCADE;
TRUNCATE TABLE vendor_items.item_images CASCADE;
TRUNCATE TABLE vendor_items.items CASCADE;
TRUNCATE TABLE vendor_items.categories CASCADE;
TRUNCATE TABLE vendor_accounts.restaurant_details CASCADE;
TRUNCATE TABLE vendor_accounts.vendor_accounts CASCADE;
TRUNCATE TABLE vendor_constants.cuisines CASCADE;
TRUNCATE TABLE user_sso.user_sessions CASCADE;
TRUNCATE TABLE user_profile.addresses CASCADE;
TRUNCATE TABLE user_profile.users CASCADE;
TRUNCATE TABLE qvickly_grocery_products.items CASCADE;
TRUNCATE TABLE qvickly_grocery_products.item_groups CASCADE;
TRUNCATE TABLE qvickly_grocery_products.dashboard_posters CASCADE;
TRUNCATE TABLE delivery_partners.delivery_partners CASCADE;

-- Insert Cuisines
INSERT INTO vendor_constants.cuisines (id, name) VALUES
                                                     ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'North Indian'),
                                                     ('b2c3d4e5-f6a7-8901-bcde-f23456789012', 'South Indian'),
                                                     ('c3d4e5f6-a7b8-9012-cdef-345678901234', 'Chinese'),
                                                     ('d4e5f6a7-b8c9-0123-def0-456789012345', 'Italian'),
                                                     ('e5f6a7b8-c9d0-1234-ef01-567890123456', 'Continental'),
                                                     ('f6a7b8c9-d0e1-2345-f012-678901234567', 'Mexican'),
                                                     ('a7b8c9d0-e1f2-3456-0123-789012345678', 'Thai'),
                                                     ('b8c9d0e1-f2a3-4567-1234-890123456789', 'Fast Food');

-- Insert Users
INSERT INTO user_profile.users (id, phone_number, email, full_name, google_id, profile_picture_url, is_marketing_opted, created_at, updated_at) VALUES
                                                                                                                                                    ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '+919876543210', 'john.doe@gmail.com', 'John Doe', 'google_123456789', 'https://example.com/avatars/john.jpg', true, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),
                                                                                                                                                    ('bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '+919876543211', 'priya.sharma@gmail.com', 'Priya Sharma', 'google_987654321', 'https://example.com/avatars/priya.jpg', true, '2024-02-20 14:45:00', '2024-02-20 14:45:00'),
                                                                                                                                                    ('cccccccc-dddd-eeee-ffff-000000000000', '+919876543212', 'raj.patel@gmail.com', 'Raj Patel', 'google_456789123', 'https://example.com/avatars/raj.jpg', false, '2024-03-10 09:15:00', '2024-03-10 09:15:00'),
                                                                                                                                                    ('dddddddd-eeee-ffff-0000-111111111111', '+919876543213', 'anita.singh@gmail.com', 'Anita Singh', 'google_789123456', 'https://example.com/avatars/anita.jpg', true, '2024-04-05 16:20:00', '2024-04-05 16:20:00'),
                                                                                                                                                    ('eeeeeeee-ffff-0000-1111-222222222222', '+919876543214', 'vikram.kumar@gmail.com', 'Vikram Kumar', 'google_321654987', 'https://example.com/avatars/vikram.jpg', true, '2024-05-01 11:10:00', '2024-05-01 11:10:00'),
                                                                                                                                                    ('ffffffff-0000-1111-2222-333333333333', '+919876543215', 'neha.gupta@gmail.com', 'Neha Gupta', 'google_654987321', 'https://example.com/avatars/neha.jpg', false, '2024-05-15 13:25:00', '2024-05-15 13:25:00');

-- Insert User Addresses
INSERT INTO user_profile.addresses (id, user_id, label, address_line1, address_line2, city, state, postal_code, country, latitude, longitude, is_default, created_at, updated_at) VALUES
                                                                                                                                                                                      ('11111111-2222-3333-4444-555555555555', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'Home', '123 MG Road', 'Near Metro Station', 'Delhi', 'Delhi', '110001', 'India', 28.613939, 77.209021, true, '2024-01-15 10:30:00', '2024-01-15 10:30:00'),
                                                                                                                                                                                      ('22222222-3333-4444-5555-666666666666', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'Office', '456 Brigade Road', 'Bangalore Central', 'Bangalore', 'Karnataka', '560001', 'India', 12.971598, 77.594562, true, '2024-02-20 14:45:00', '2024-02-20 14:45:00'),
                                                                                                                                                                                      ('33333333-4444-5555-6666-777777777777', 'cccccccc-dddd-eeee-ffff-000000000000', 'Home', '789 Marine Drive', 'Nariman Point', 'Mumbai', 'Maharashtra', '400001', 'India', 18.922881, 72.834632, true, '2024-03-10 09:15:00', '2024-03-10 09:15:00'),
                                                                                                                                                                                      ('44444444-5555-6666-7777-888888888888', 'dddddddd-eeee-ffff-0000-111111111111', 'Home', '321 Park Street', 'Near Victoria Memorial', 'Kolkata', 'West Bengal', '700016', 'India', 22.544249, 88.342182, true, '2024-04-05 16:20:00', '2024-04-05 16:20:00'),
                                                                                                                                                                                      ('55555555-6666-7777-8888-999999999999', 'eeeeeeee-ffff-0000-1111-222222222222', 'Work', '654 Anna Salai', 'Mount Road', 'Chennai', 'Tamil Nadu', '600002', 'India', 13.061415, 80.249914, true, '2024-05-01 11:10:00', '2024-05-01 11:10:00'),
                                                                                                                                                                                      ('66666666-7777-8888-9999-aaaaaaaaaaaa', 'ffffffff-0000-1111-2222-333333333333', 'Home', '987 Banjara Hills', 'Road No 12', 'Hyderabad', 'Telangana', '500034', 'India', 17.412563, 78.448639, true, '2024-05-15 13:25:00', '2024-05-15 13:25:00');

-- Insert User Sessions
INSERT INTO user_sso.user_sessions (id, user_id, refresh_token, refresh_token_expires_at, ip_address, user_agent, device_info, created_at, last_seen_at) VALUES
                                                                                                                                                             ('12345678-abcd-ef01-2345-6789abcdef01', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'refresh_token_abc123', '2024-07-15 10:30:00', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', '{"device": "desktop", "os": "Windows", "browser": "Chrome"}', '2024-05-15 10:30:00', '2024-06-01 15:45:00'),
                                                                                                                                                             ('23456789-bcde-f012-3456-789abcdef012', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'refresh_token_def456', '2024-07-20 14:45:00', '192.168.1.101', 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)', '{"device": "mobile", "os": "iOS", "browser": "Safari"}', '2024-05-20 14:45:00', '2024-06-02 09:30:00'),
                                                                                                                                                             ('34567890-cdef-0123-4567-89abcdef0123', 'cccccccc-dddd-eeee-ffff-000000000000', 'refresh_token_ghi789', '2024-07-10 09:15:00', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36', '{"device": "desktop", "os": "macOS", "browser": "Chrome"}', '2024-05-10 09:15:00', '2024-06-03 12:20:00');

-- Insert Vendor Accounts
INSERT INTO vendor_accounts.vendor_accounts (id, phone_number, account_type, business_name, owner_name, email, address, created_at, latitude, longitude, gstin_number, opening_time, closing_time, image_url, live_status) VALUES
                                                                                                                                                                                                                               ('69fc8ba4-d11f-4618-9042-c1523d381013', '8010201921', 'restaurant', 'Spice Garden Restaurant', 'Rajat Sharma', 'rajatnd9@gmail.com', '123 Commercial Street, Koramangala', '2024-01-10 08:00:00', 12.9352, 77.6245, 'GST123456789ABC', '09:00:00', '23:00:00', 'https://example.com/restaurants/spice-garden.jpg', true),
                                                                                                                                                                                                                               ('77777777-8888-9999-aaaa-bbbbbbbbbbbb', '9876543210', 'restaurant', 'Delhi Darbar', 'Amit Kumar', 'amit.delhi@gmail.com', '456 CP Market, Connaught Place', '2024-01-15 10:00:00', 28.6315, 77.2167, 'GST987654321DEF', '10:00:00', '22:30:00', 'https://example.com/restaurants/delhi-darbar.jpg', true),
                                                                                                                                                                                                                               ('88888888-9999-aaaa-bbbb-cccccccccccc', '9876543211', 'store', 'Fresh Mart Grocery', 'Priya Patel', 'priya.freshmart@gmail.com', '789 Link Road, Bandra West', '2024-02-01 07:30:00', 19.0596, 72.8295, 'GST456789123GHI', '07:00:00', '22:00:00', 'https://example.com/stores/fresh-mart.jpg', true),
                                                                                                                                                                                                                               ('99999999-aaaa-bbbb-cccc-dddddddddddd', '9876543212', 'restaurant', 'Pizza Corner', 'Ravi Singh', 'ravi.pizza@gmail.com', '321 Brigade Road, MG Road', '2024-02-10 11:00:00', 12.9716, 77.5946, 'GST789123456JKL', '11:00:00', '23:30:00', 'https://example.com/restaurants/pizza-corner.jpg', true),
                                                                                                                                                                                                                               ('aaaabbbb-cccc-dddd-eeee-ffffffffffff', '9876543213', 'store', 'Green Valley Organics', 'Sunita Reddy', 'sunita.organics@gmail.com', '654 Jubilee Hills, Road No 36', '2024-02-20 06:00:00', 17.4065, 78.4772, 'GST321654987MNO', '06:00:00', '21:00:00', 'https://example.com/stores/green-valley.jpg', true);

-- Insert Restaurant Details
INSERT INTO vendor_accounts.restaurant_details (account_id, fssai_license_no, cuisine_id) VALUES
                                                                                              ('69fc8ba4-d11f-4618-9042-c1523d381013', 'FSSAI12345678901', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890'),
                                                                                              ('77777777-8888-9999-aaaa-bbbbbbbbbbbb', 'FSSAI98765432109', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890'),
                                                                                              ('99999999-aaaa-bbbb-cccc-dddddddddddd', 'FSSAI45678912345', 'd4e5f6a7-b8c9-0123-def0-456789012345');

-- Insert Delivery Partners
INSERT INTO delivery_partners.delivery_partners (id, name, phone_number, pin, is_active, created_at, updated_at) VALUES
                                                                                                                     ('de111111-2222-3333-4444-555555555555', 'Ramesh Kumar', '9988776655', '1234', true, '2024-01-05 08:00:00', '2024-01-05 08:00:00'),
                                                                                                                     ('de222222-3333-4444-5555-666666666666', 'Suresh Patel', '9988776656', '5678', true, '2024-01-10 09:00:00', '2024-01-10 09:00:00'),
                                                                                                                     ('de333333-4444-5555-6666-777777777777', 'Mahesh Singh', '9988776657', '9012', true, '2024-01-15 10:00:00', '2024-01-15 10:00:00'),
                                                                                                                     ('de444444-5555-6666-7777-888888888888', 'Ganesh Reddy', '9988776658', '3456', false, '2024-02-01 11:00:00', '2024-02-01 11:00:00');

-- Insert Item Groups for Qvickly Grocery Products
INSERT INTO qvickly_grocery_products.item_groups (id, name, slug, display_order) VALUES
                                                                                     (1, 'Fruits & Vegetables', 'fruits-vegetables', 1),
                                                                                     (2, 'Dairy & Eggs', 'dairy-eggs', 2),
                                                                                     (3, 'Bakery & Snacks', 'bakery-snacks', 3),
                                                                                     (4, 'Beverages', 'beverages', 4),
                                                                                     (5, 'Personal Care', 'personal-care', 5),
                                                                                     (6, 'Household Items', 'household-items', 6);

-- Insert Qvickly Grocery Items
INSERT INTO qvickly_grocery_products.items (id, name, image_url, group_id, price, is_available, created_at, rating, description, sustainability_facts) VALUES
                                                                                                                                                           (1, 'Fresh Bananas (1kg)', 'https://example.com/products/bananas.jpg', 1, 40.00, true, '2024-01-01 00:00:00', 4.5, 'Fresh ripe bananas rich in potassium', 'Locally sourced, minimal packaging'),
                                                                                                                                                           (2, 'Red Apples (1kg)', 'https://example.com/products/apples.jpg', 1, 120.00, true, '2024-01-01 00:00:00', 4.3, 'Crisp red apples from Kashmir', 'Organic farming practices'),
                                                                                                                                                           (3, 'Fresh Milk (1L)', 'https://example.com/products/milk.jpg', 2, 55.00, true, '2024-01-01 00:00:00', 4.7, 'Pure cow milk from local dairy', 'Glass bottles for recycling'),
                                                                                                                                                           (4, 'Brown Bread', 'https://example.com/products/bread.jpg', 3, 35.00, true, '2024-01-01 00:00:00', 4.2, 'Whole wheat brown bread', 'Minimal preservatives'),
                                                                                                                                                           (5, 'Orange Juice (1L)', 'https://example.com/products/orange-juice.jpg', 4, 85.00, true, '2024-01-01 00:00:00', 4.4, 'Fresh orange juice with pulp', 'No artificial colors or flavors'),
                                                                                                                                                           (6, 'Shampoo (200ml)', 'https://example.com/products/shampoo.jpg', 5, 150.00, true, '2024-01-01 00:00:00', 4.1, 'Herbal shampoo for all hair types', 'Biodegradable formula');

-- Insert Dashboard Posters
INSERT INTO qvickly_grocery_products.dashboard_posters (id, image_url, link_url, active, created_at) VALUES
                                                                                                         (1, 'https://example.com/posters/summer-sale.jpg', 'https://example.com/offers/summer-sale', true, '2024-05-01 00:00:00'),
                                                                                                         (2, 'https://example.com/posters/fresh-fruits.jpg', 'https://example.com/categories/fruits', true, '2024-05-15 00:00:00'),
                                                                                                         (3, 'https://example.com/posters/dairy-special.jpg', 'https://example.com/categories/dairy', false, '2024-04-01 00:00:00');

-- Insert Vendor Item Categories
INSERT INTO vendor_items.categories (id, name, created_at) VALUES
                                                               (1, 'Starters', '2024-01-01 00:00:00'),
                                                               (2, 'Main Course', '2024-01-01 00:00:00'),
                                                               (3, 'Desserts', '2024-01-01 00:00:00'),
                                                               (4, 'Beverages', '2024-01-01 00:00:00'),
                                                               (5, 'Rice & Biryani', '2024-01-01 00:00:00'),
                                                               (6, 'Pizza', '2024-01-01 00:00:00'),
                                                               (7, 'Groceries', '2024-01-01 00:00:00');

-- Insert Vendor Items
INSERT INTO vendor_items.items (id, account_id, category_id, name, description, price_retail, price_wholesale, is_available, stock, created_at, updated_at, search_keywords, is_active, vendor_id) VALUES
                                                                                                                                                                                                       (1, '69fc8ba4-d11f-4618-9042-c1523d381013', 1, 'Chicken Tikka', 'Grilled chicken marinated in spices', 180.00, 150.00, true, 50, '2024-01-10 08:00:00', '2024-01-10 08:00:00', 'chicken tikka starter grilled', true, '69fc8ba4-d11f-4618-9042-c1523d381013'),
                                                                                                                                                                                                       (2, '69fc8ba4-d11f-4618-9042-c1523d381013', 2, 'Butter Chicken', 'Creamy tomato-based chicken curry', 320.00, 280.00, true, 30, '2024-01-10 08:00:00', '2024-01-10 08:00:00', 'butter chicken curry main course', true, '69fc8ba4-d11f-4618-9042-c1523d381013'),
                                                                                                                                                                                                       (3, '69fc8ba4-d11f-4618-9042-c1523d381013', 5, 'Chicken Biryani', 'Aromatic basmati rice with spiced chicken', 280.00, 240.00, true, 25, '2024-01-10 08:00:00', '2024-01-10 08:00:00', 'chicken biryani rice basmati', true, '69fc8ba4-d11f-4618-9042-c1523d381013'),
                                                                                                                                                                                                       (4, '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 1, 'Paneer Tikka', 'Grilled cottage cheese with vegetables', 160.00, 130.00, true, 40, '2024-01-15 10:00:00', '2024-01-15 10:00:00', 'paneer tikka vegetarian starter', true, '77777777-8888-9999-aaaa-bbbbbbbbbbbb'),
                                                                                                                                                                                                       (5, '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 2, 'Dal Makhani', 'Rich black lentil curry', 220.00, 180.00, true, 35, '2024-01-15 10:00:00', '2024-01-15 10:00:00', 'dal makhani lentil curry', true, '77777777-8888-9999-aaaa-bbbbbbbbbbbb'),
                                                                                                                                                                                                       (6, '99999999-aaaa-bbbb-cccc-dddddddddddd', 6, 'Margherita Pizza', 'Classic pizza with tomato and mozzarella', 250.00, 200.00, true, 20, '2024-02-10 11:00:00', '2024-02-10 11:00:00', 'margherita pizza cheese tomato', true, '99999999-aaaa-bbbb-cccc-dddddddddddd'),
                                                                                                                                                                                                       (7, '99999999-aaaa-bbbb-cccc-dddddddddddd', 6, 'Pepperoni Pizza', 'Pizza topped with pepperoni and cheese', 350.00, 300.00, true, 15, '2024-02-10 11:00:00', '2024-02-10 11:00:00', 'pepperoni pizza meat cheese', true, '99999999-aaaa-bbbb-cccc-dddddddddddd'),
                                                                                                                                                                                                       (8, '88888888-9999-aaaa-bbbb-cccccccccccc', 7, 'Organic Tomatoes (1kg)', 'Fresh organic tomatoes', 60.00, 45.00, true, 100, '2024-02-01 07:30:00', '2024-02-01 07:30:00', 'organic tomatoes vegetables fresh', true, '88888888-9999-aaaa-bbbb-cccccccccccc'),
                                                                                                                                                                                                       (9, '88888888-9999-aaaa-bbbb-cccccccccccc', 7, 'Basmati Rice (5kg)', 'Premium quality basmati rice', 450.00, 400.00, true, 50, '2024-02-01 07:30:00', '2024-02-01 07:30:00', 'basmati rice premium grain', true, '88888888-9999-aaaa-bbbb-cccccccccccc'),
                                                                                                                                                                                                       (10, 'aaaabbbb-cccc-dddd-eeee-ffffffffffff', 7, 'Organic Spinach (500g)', 'Fresh organic spinach leaves', 35.00, 28.00, true, 75, '2024-02-20 06:00:00', '2024-02-20 06:00:00', 'organic spinach green leafy vegetables', true, 'aaaabbbb-cccc-dddd-eeee-ffffffffffff');

-- Insert Item Images
INSERT INTO vendor_items.item_images (id, item_id, image_url, "position", created_at) VALUES
                                                                                          (1, 1, 'https://example.com/items/chicken-tikka-1.jpg', 1, '2024-01-10 08:00:00'),
                                                                                          (2, 1, 'https://example.com/items/chicken-tikka-2.jpg', 2, '2024-01-10 08:00:00'),
                                                                                          (3, 2, 'https://example.com/items/butter-chicken-1.jpg', 1, '2024-01-10 08:00:00'),
                                                                                          (4, 3, 'https://example.com/items/chicken-biryani-1.jpg', 1, '2024-01-10 08:00:00'),
                                                                                          (5, 4, 'https://example.com/items/paneer-tikka-1.jpg', 1, '2024-01-15 10:00:00'),
                                                                                          (6, 5, 'https://example.com/items/dal-makhani-1.jpg', 1, '2024-01-15 10:00:00'),
                                                                                          (7, 6, 'https://example.com/items/margherita-pizza-1.jpg', 1, '2024-02-10 11:00:00'),
                                                                                          (8, 7, 'https://example.com/items/pepperoni-pizza-1.jpg', 1, '2024-02-10 11:00:00'),
                                                                                          (9, 8, 'https://example.com/items/organic-tomatoes-1.jpg', 1, '2024-02-01 07:30:00'),
                                                                                          (10, 9, 'https://example.com/items/basmati-rice-1.jpg', 1, '2024-02-01 07:30:00'),
                                                                                          (11, 10, 'https://example.com/items/organic-spinach-1.jpg', 1, '2024-02-20 06:00:00');

-- Insert Vendor Inventory
INSERT INTO public.vendor_inventory (id, vendor_id, item_id, stock_quantity, is_available, price_override, created_at, updated_at) VALUES
                                                                                                                                       ('aa111111-2222-3333-4444-555555555555', '69fc8ba4-d11f-4618-9042-c1523d381013', 1, 50, true, NULL, '2024-01-10 08:00:00', '2024-01-10 08:00:00'),
                                                                                                                                       ('bb222222-3333-4444-5555-666666666666', '69fc8ba4-d11f-4618-9042-c1523d381013', 2, 30, true, NULL, '2024-01-10 08:00:00', '2024-01-10 08:00:00'),
                                                                                                                                       ('cc333333-4444-5555-6666-777777777777', '69fc8ba4-d11f-4618-9042-c1523d381013', 3, 25, true, NULL, '2024-01-10 08:00:00', '2024-01-10 08:00:00'),
                                                                                                                                       ('dd444444-5555-6666-7777-888888888888', '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 4, 40, true, 150.00, '2024-01-15 10:00:00', '2024-01-15 10:00:00'),
                                                                                                                                       ('ee555555-6666-7777-8888-999999999999', '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 5, 35, true, NULL, '2024-01-15 10:00:00', '2024-01-15 10:00:00'),
                                                                                                                                       ('ff666666-7777-8888-9999-aaaaaaaaaaaa', '99999999-aaaa-bbbb-cccc-dddddddddddd', 6, 20, true, NULL, '2024-02-10 11:00:00', '2024-02-10 11:00:00'),
                                                                                                                                       ('aa777777-8888-9999-aaaa-bbbbbbbbbbbb', '99999999-aaaa-bbbb-cccc-dddddddddddd', 7, 15, true, NULL, '2024-02-10 11:00:00', '2024-02-10 11:00:00'),
                                                                                                                                       ('bb888888-9999-aaaa-bbbb-cccccccccccc', '88888888-9999-aaaa-bbbb-cccccccccccc', 8, 100, true, 55.00, '2024-02-01 07:30:00', '2024-02-01 07:30:00'),
                                                                                                                                       ('cc999999-aaaa-bbbb-cccc-dddddddddddd', '88888888-9999-aaaa-bbbb-cccccccccccc', 9, 50, true, NULL, '2024-02-01 07:30:00', '2024-02-01 07:30:00'),
                                                                                                                                       ('ddaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'aaaabbbb-cccc-dddd-eeee-ffffffffffff', 10, 75, true, NULL, '2024-02-20 06:00:00', '2024-02-20 06:00:00');

-- Insert Inventory Movements
INSERT INTO public.inventory_movements (id, vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_by, created_at) VALUES
                                                                                                                                                                      ('a1111111-2222-3333-4444-555555555555', 'aa111111-2222-3333-4444-555555555555', 'stock_in', 50, 0, 50, 'Initial stock', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '2024-01-10 08:00:00'),
                                                                                                                                                                      ('b2222222-3333-4444-5555-666666666666', 'bb222222-3333-4444-5555-666666666666', 'stock_in', 30, 0, 30, 'Initial stock', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '2024-01-10 08:00:00'),
                                                                                                                                                                      ('c3333333-4444-5555-6666-777777777777', 'aa111111-2222-3333-4444-555555555555', 'sale', -2, 50, 48, 'Order fulfillment', NULL, '2024-01-15 12:30:00'),
                                                                                                                                                                      ('d4444444-5555-6666-7777-888888888888', 'bb888888-9999-aaaa-bbbb-cccccccccccc', 'stock_in', 100, 0, 100, 'Fresh produce delivery', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '2024-02-01 07:30:00'),
                                                                                                                                                                      ('e5555555-6666-7777-8888-999999999999', 'bb888888-9999-aaaa-bbbb-cccccccccccc', 'adjustment', -5, 100, 95, 'Damaged items removed', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '2024-02-02 10:00:00'),
                                                                                                                                                                      ('f6666666-7777-8888-9999-aaaaaaaaaaaa', 'cc999999-aaaa-bbbb-cccc-dddddddddddd', 'stock_in', 50, 0, 50, 'Initial rice stock', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '2024-02-01 07:30:00'),
                                                                                                                                                                      ('a7777777-8888-9999-aaaa-bbbbbbbbbbbb', 'ddaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'stock_in', 75, 0, 75, 'Fresh spinach delivery', 'cccccccc-dddd-eeee-ffff-000000000000', '2024-02-20 06:00:00'),
                                                                                                                                                                      ('b8888888-9999-aaaa-bbbb-cccccccccccc', 'ff666666-7777-8888-9999-aaaaaaaaaaaa', 'sale', -3, 20, 17, 'Pizza orders', NULL, '2024-02-15 19:30:00');

-- Insert Orders
INSERT INTO orders.orders (id, account_id, user_id, order_time, total_amount, status, location, created_at, updated_at, customer_id, delivery_partner_id, customer_name, customer_address_id, pack_by_time, paid_by_time, delivered_by_time) VALUES
                                                                                                                                                                                                                                                 (1, '69fc8ba4-d11f-4618-9042-c1523d381013', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '2024-05-22 12:30:00', 500.00, 'completed', 'Delhi', '2024-05-22 12:30:00', '2024-05-22 14:45:00', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'de111111-2222-3333-4444-555555555555', 'John Doe', '11111111-2222-3333-4444-555555555555', '2024-05-22 13:00:00', '2024-05-22 12:35:00', '2024-05-22 14:45:00'),
                                                                                                                                                                                                                                                 (2, '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '2024-05-23 19:15:00', 380.00, 'ready', 'Bangalore', '2024-05-23 19:15:00', '2024-05-23 20:30:00', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'de222222-3333-4444-5555-666666666666', 'Priya Sharma', '22222222-3333-4444-5555-666666666666', '2024-05-23 20:45:00', '2024-05-23 19:20:00', NULL),
                                                                                                                                                                                                                                                 (3, '99999999-aaaa-bbbb-cccc-dddddddddddd', 'cccccccc-dddd-eeee-ffff-000000000000', '2024-05-24 20:45:00', 600.00, 'packed', 'Mumbai', '2024-05-24 20:45:00', '2024-05-24 21:15:00', 'cccccccc-dddd-eeee-ffff-000000000000', NULL, 'Raj Patel', '33333333-4444-5555-6666-777777777777', '2024-05-24 21:30:00', '2024-05-24 20:50:00', NULL),
                                                                                                                                                                                                                                                 (4, '69fc8ba4-d11f-4618-9042-c1523d381013', 'dddddddd-eeee-ffff-0000-111111111111', '2024-05-25 13:20:00', 280.00, 'accepted', 'Kolkata', '2024-05-25 13:20:00', '2024-05-25 13:45:00', 'dddddddd-eeee-ffff-0000-111111111111', NULL, 'Anita Singh', '44444444-5555-6666-7777-888888888888', '2024-05-25 14:30:00', '2024-05-25 13:25:00', NULL),
                                                                                                                                                                                                                                                 (5, '88888888-9999-aaaa-bbbb-cccccccccccc', 'eeeeeeee-ffff-0000-1111-222222222222', '2024-05-26 10:10:00', 510.00, 'pending', 'Chennai', '2024-05-26 10:10:00', '2024-05-26 10:10:00', 'eeeeeeee-ffff-0000-1111-222222222222', NULL, 'Vikram Kumar', '55555555-6666-7777-8888-999999999999', NULL, NULL, NULL),
                                                                                                                                                                                                                                                 (6, '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 'ffffffff-0000-1111-2222-333333333333', '2024-05-20 18:30:00', 220.00, 'rejected', 'Hyderabad', '2024-05-20 18:30:00', '2024-05-20 19:00:00', 'ffffffff-0000-1111-2222-333333333333', NULL, 'Neha Gupta', '66666666-7777-8888-9999-aaaaaaaaaaaa', NULL, NULL, NULL),
                                                                                                                                                                                                                                                 (7, 'aaaabbbb-cccc-dddd-eeee-ffffffffffff', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '2024-05-27 08:15:00', 70.00, 'pending', 'Delhi', '2024-05-27 08:15:00', '2024-05-27 08:15:00', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', NULL, 'John Doe', '11111111-2222-3333-4444-555555555555', NULL, NULL, NULL),
                                                                                                                                                                                                                                                 (8, '99999999-aaaa-bbbb-cccc-dddddddddddd', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '2024-05-28 14:20:00', 750.00, 'accepted', 'Bangalore', '2024-05-28 14:20:00', '2024-05-28 14:35:00', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'de333333-4444-5555-6666-777777777777', 'Priya Sharma', '22222222-3333-4444-5555-666666666666', '2024-05-28 15:30:00', '2024-05-28 14:25:00', NULL);

-- Insert Order Items
INSERT INTO orders.order_items (id, order_id, item_id, quantity, instructions, is_checked) VALUES
                                                                                               (1, 1, 1, 2, 'Medium spice level', true),
                                                                                               (2, 1, 2, 1, 'Extra creamy', true),
                                                                                               (3, 1, 3, 1, 'Less oil please', true),
                                                                                               (4, 2, 4, 2, 'Well grilled', true),
                                                                                               (5, 2, 5, 1, 'Extra butter', false),
                                                                                               (6, 3, 6, 1, 'Thin crust', true),
                                                                                               (7, 3, 7, 1, 'Extra cheese', false),
                                                                                               (8, 4, 3, 1, 'Fragrant rice', false),
                                                                                               (9, 5, 8, 3, 'Fresh and ripe', false),
                                                                                               (10, 5, 9, 1, 'Premium quality', false),
                                                                                               (11, 6, 5, 1, 'Regular preparation', false),
                                                                                               (12, 7, 10, 2, 'Fresh green leaves', false),
                                                                                               (13, 8, 6, 2, 'Make it crispy', false),
                                                                                               (14, 8, 7, 1, 'Extra pepperoni', false);

-- Insert Order Status Logs
INSERT INTO orders.order_status_logs (id, order_id, status, changed_at) VALUES
                                                                            (1, 1, 'pending', '2024-05-22 12:30:00'),
                                                                            (2, 1, 'accepted', '2024-05-22 12:35:00'),
                                                                            (3, 1, 'packed', '2024-05-22 13:00:00'),
                                                                            (4, 1, 'ready', '2024-05-22 13:15:00'),
                                                                            (5, 1, 'completed', '2024-05-22 14:45:00'),
                                                                            (6, 2, 'pending', '2024-05-23 19:15:00'),
                                                                            (7, 2, 'accepted', '2024-05-23 19:20:00'),
                                                                            (8, 2, 'packed', '2024-05-23 20:00:00'),
                                                                            (9, 2, 'ready', '2024-05-23 20:30:00'),
                                                                            (10, 3, 'pending', '2024-05-24 20:45:00'),
                                                                            (11, 3, 'accepted', '2024-05-24 20:50:00'),
                                                                            (12, 3, 'packed', '2024-05-24 21:15:00'),
                                                                            (13, 4, 'pending', '2024-05-25 13:20:00'),
                                                                            (14, 4, 'accepted', '2024-05-25 13:45:00'),
                                                                            (15, 5, 'pending', '2024-05-26 10:10:00'),
                                                                            (16, 6, 'pending', '2024-05-20 18:30:00'),
                                                                            (17, 6, 'rejected', '2024-05-20 19:00:00'),
                                                                            (18, 7, 'pending', '2024-05-27 08:15:00'),
                                                                            (19, 8, 'pending', '2024-05-28 14:20:00'),
                                                                            (20, 8, 'accepted', '2024-05-28 14:35:00');

