CREATE TABLE IF NOT EXISTS user
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name       text                        NOT NULL,
    password text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS products
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text                        NOT NULL,
    description     text,
    price int
);
create table if not exists shop
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title           text                        NOT NULL,
    description     text
);
CREATE TABLE IF NOT EXISTS category_and_products
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    "updated_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    "category" bigserial,
    "products"       bigserial,
    FOREIGN KEY (category)
        REFERENCES products(id),
    FOREIGN KEY (products)
        REFERENCES category(id)
);