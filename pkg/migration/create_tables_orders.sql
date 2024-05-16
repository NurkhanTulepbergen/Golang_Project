CREATE TABLE IF NOT EXISTS order_helper
(
    product_id      int,
    product_name    varchar(255),
    product_price   float NOT NULL,
    quantity        int NOT NULL DEFAULT 1,
    order_id        int not null ,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE OR REPLACE FUNCTION update_product_price()
    RETURNS TRIGGER AS $$
BEGIN
    -- Заполняем поле product_price из таблицы products
    SELECT price INTO NEW.product_price FROM products WHERE id = NEW.product_id;
    -- Заполняем поле product_name из таблицы products
    SELECT title INTO NEW.product_name FROM products WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_order_helper
    BEFORE INSERT ON order_helper
    FOR EACH ROW
EXECUTE FUNCTION update_product_price();

/*EXAMPLE OF INSERTING
INSERT INTO order_helper (product_id, quantity, order_id)
VALUES (1, 2, 1);*/
----------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS orders
(
    id              bigserial PRIMARY KEY,
    user_id         int NOT NULL,
    products        jsonb,
    total_amount    float NOT NULL,
    delivery_addr   text NOT NULL,
    status          varchar(255),
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE OR REPLACE FUNCTION update_order_total()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET total_amount = (
        SELECT SUM(product_price * quantity)
        FROM order_helper
        WHERE order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_or_update_order_helper
    BEFORE INSERT OR UPDATE ON order_helper
    FOR EACH ROW
EXECUTE FUNCTION update_order_total();

drop table orders;