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
    price           float
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

CREATE TABLE IF NOT EXISTS tokens
(
    hash    BYTEA PRIMARY KEY,
    user_id BIGINT                      NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry  TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope   TEXT                        NOT NULL
);

CREATE TABLE Orders (
                        ID VARCHAR(255) PRIMARY KEY,
                        "user" VARCHAR(255),
                        TotalAmount NUMERIC(10, 2),
                        DeliveryAddr VARCHAR(255),
                        Status VARCHAR(50),
                        CreatedAt TIMESTAMP
);

CREATE TABLE Cart (
                      UserID VARCHAR(255) PRIMARY KEY,
                      Items JSONB
);

CREATE TABLE IF NOT EXISTS permissions (
                                           id bigserial PRIMARY KEY,
                                           code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
                                                 user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                                 permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
                                                 PRIMARY KEY (user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
    ('shop:read'),
    ('shop:write'),
    ('catalog:read'),
    ('catalog:write');


------------------------------------------------------------------------------------------------------------------------


INSERT INTO users (name, email, password, activated)
VALUES ('Nurkhan Tulepbergen', 'n_tulepbegren@kbtu.kz', 'Nurkhan05', true);

------------------------------------------------------------------------------------------------------------------------

INSERT INTO shop (title, description, type)
VALUES ('Book World', 'Here is the shop where you can find every book for your self', 'Book');


INSERT INTO products (title, description, price)
VALUES ('Harry Potter and the Philosophers Stone', 'The first book in the Harry Potter series', 15.99);

INSERT INTO products (title, description, price)
VALUES ('To Kill a Mockingbird', 'A classic novel by Harper Lee', 12.99);

INSERT INTO products (title, description, price)
VALUES ('The Great Gatsby', 'A novel by F. Scott Fitzgerald', 11.99);

INSERT INTO products (title, description, price)
VALUES ('1984', 'A dystopian novel by George Orwell', 14.99);

INSERT INTO products (title, description, price)
VALUES ('The Catcher in the Rye', 'A novel by J.D. Salinger', 10.99);


INSERT INTO shop_and_products (shop_id, product_id)
VALUES (1, 1); -- Linking 'Harry Potter and the Philosopher's Stone' to 'Book World'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (1, 2); -- Linking 'To Kill a Mockingbird' to 'Book World'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (1, 3); -- Linking 'The Great Gatsby' to 'Book World'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (1, 4); -- Linking '1984' to 'Book World'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (1, 5); -- Linking 'The Catcher in the Rye' to 'Book World'

------------------------------------------------------------------------------------------------------------------------

INSERT INTO shop (title, description, type)
VALUES ('iSpace', 'All apple devices in one shop', 'Electronic');


INSERT INTO products (title, description, price)
VALUES ('iPhone 13 Pro', 'Apples latest flagship smartphone', 999.99);

INSERT INTO products (title, description, price)
VALUES ('MacBook Pro', 'Powerful laptop by Apple', 1499.99);

INSERT INTO products (title, description, price)
VALUES ('iPad Air', 'Versatile tablet by Apple', 599.99);

INSERT INTO products (title, description, price)
VALUES ('AirPods Pro', 'High-quality wireless earbuds by Apple', 249.99);

INSERT INTO products (title, description, price)
VALUES ('Apple Watch Series 7', 'Advanced smartwatch with health features', 399.99);


INSERT INTO shop_and_products (shop_id, product_id)
VALUES (2, 6); -- Linking 'iPhone 13 Pro' to 'iSpace'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (2, 7); -- Linking 'MacBook Pro' to 'iSpace'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (2, 8); -- Linking 'iPad Air' to 'iSpace'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (2, 9); -- Linking 'AirPods Pro' to 'iSpace'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (2, 10); -- Linking 'Apple Watch Series 7' to 'iSpace'

------------------------------------------------------------------------------------------------------------------------

INSERT INTO shop (title, description, type)
VALUES ('Samsung', 'All samsung devices in one shop', 'Electronic');


INSERT INTO products (title, description, price)
VALUES ('Samsung Galaxy S21 Ultra', 'Flagship smartphone with powerful camera', 1199.99);

INSERT INTO products (title, description, price)
VALUES ('Samsung Galaxy Tab S7', 'High-performance Android tablet', 649.99);

INSERT INTO products (title, description, price)
VALUES ('Samsung Galaxy Watch 4', 'Feature-rich smartwatch with health tracking', 299.99);

INSERT INTO products (title, description, price)
VALUES ('Samsung Galaxy Buds Pro', 'Premium wireless earbuds with active noise cancellation', 199.99);

INSERT INTO products (title, description, price)
VALUES ('Samsung Odyssey G7', 'Curved gaming monitor with QLED technology', 699.99);


INSERT INTO shop_and_products (shop_id, product_id)
VALUES (3, 11); -- Linking 'Samsung Galaxy S21 Ultra' to 'Samsung'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (3, 12); -- Linking 'Samsung Galaxy Tab S7' to 'Samsung'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (3, 13); -- Linking 'Samsung Galaxy Watch 4' to 'Samsung'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (3, 14); -- Linking 'Samsung Galaxy Buds Pro' to 'Samsung'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (3, 15); -- Linking 'Samsung Odyssey G7' to 'Samsung'

------------------------------------------------------------------------------------------------------------------------

INSERT INTO shop (title, description, type)
VALUES ('DrinkLand', 'Every drink in your mind is here', 'Drink');


INSERT INTO products (title, description, price)
VALUES ('Coca-Cola', 'Refreshing carbonated soft drink', 1.99);

INSERT INTO products (title, description, price)
VALUES ('Pepsi', 'Another popular carbonated soft drink', 1.79);

INSERT INTO products (title, description, price)
VALUES ('Red Bull', 'Energy drink for increased performance', 2.49);

INSERT INTO products (title, description, price)
VALUES ('Fanta', 'Fruity carbonated soft drink', 1.69);

INSERT INTO products (title, description, price)
VALUES ('Gatorade', 'Sports drink for hydration and replenishment', 2.29);


INSERT INTO shop_and_products (shop_id, product_id)
VALUES (4, 16); -- Linking 'Coca-Cola' to 'DrinkLand'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (4, 17); -- Linking 'Pepsi' to 'DrinkLand'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (4, 18); -- Linking 'Red Bull' to 'DrinkLand'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (4, 19); -- Linking 'Fanta' to 'DrinkLand'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (4, 20); -- Linking 'Gatorade' to 'DrinkLand'

------------------------------------------------------------------------------------------------------------------------

INSERT INTO shop (title, description, type)
VALUES ('HomeComfort', 'All furniture that you imagine for your home design is here', 'Furniture');


INSERT INTO products (title, description, price)
VALUES ('Sofa', 'Comfortable seating for your living room', 699.99);

INSERT INTO products (title, description, price)
VALUES ('Dining Table', 'Sturdy table for family meals', 499.99);

INSERT INTO products (title, description, price)
VALUES ('Bed Frame', 'Elegant bed frame for a good nights sleep', 799.99);

INSERT INTO products (title, description, price)
VALUES ('Wardrobe', 'Spacious wardrobe for organizing clothes', 599.99);

INSERT INTO products (title, description, price)
VALUES ('Coffee Table', 'Functional table for your living room', 299.99);


INSERT INTO shop_and_products (shop_id, product_id)
VALUES (5, 21); -- Linking 'Sofa' to 'HomeComfort'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (5, 22); -- Linking 'Dining Table' to 'HomeComfort'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (5, 23); -- Linking 'Bed Frame' to 'HomeComfort'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (5, 24); -- Linking 'Wardrobe' to 'HomeComfort'

INSERT INTO shop_and_products (shop_id, product_id)
VALUES (5, 25); -- Linking 'Coffee Table' to 'HomeComfort'

------------------------------------------------------------------------------------------------------------------------
