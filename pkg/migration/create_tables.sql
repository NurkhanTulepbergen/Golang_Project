-- add citext extension
CREATE EXTENSION IF NOT EXISTS citext;

-- Admin@kbtu.kz = admin@kbtu.kz
CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name          TEXT                        NOT NULL,
    email         CITEXT UNIQUE               NOT NULL,
    password BYTEA                       NOT NULL,
    activated     BOOL                        NOT NULL,
    version       INTEGER                     NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS products
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text NOT NULL,
    description     text,
    price           float,
    shop_id         bigint NOT NULL,
    FOREIGN KEY (shop_id) REFERENCES shop (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS shop
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text NOT NULL,
    description     text,
    type text
);

CREATE TABLE IF NOT EXISTS shop_and_products
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    shop_id         bigint,
    product_id      bigint,
    FOREIGN KEY (shop_id) REFERENCES shop(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
drop table products;