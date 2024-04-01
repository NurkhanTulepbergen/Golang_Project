CREATE TABLE IF NOT EXISTS users
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name        text NOT NULL,
    password    text NOT NULL
    );

CREATE TABLE IF NOT EXISTS products
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text NOT NULL,
    description     text,
    price           float
    );

CREATE TABLE IF NOT EXISTS shop
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text NOT NULL,
    description     text
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
