create table if not exists follow_list
(
    user_id int,
    product_id int,
    product_name varchar(255),
    product_description text,
    product_price float,
    foreign key(user_id) references users(id),
    foreign key(product_id) references products(id)
);

CREATE OR REPLACE FUNCTION update_product_information()
    RETURNS TRIGGER AS $$
BEGIN
    -- Заполняем поле product_name из таблицы products
    SELECT title INTO NEW.product_name FROM products WHERE id = NEW.product_id;

    -- Заполняем поле product_price из таблицы products
    SELECT description INTO NEW.product_description FROM products WHERE id = NEW.product_id;

    -- Заполняем поле product_price из таблицы products
    SELECT price INTO NEW.product_price FROM products WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_follow_list
    BEFORE INSERT ON follow_list
    FOR EACH ROW
EXECUTE FUNCTION update_product_information();


create table if not exists follow
(
    user_id int,
    user_name varchar(255),
    followed_products jsonb,
    foreign key(user_id) references users(id)
);

ALTER TABLE follow
    ADD CONSTRAINT unique_user_id
        UNIQUE (user_id);

CREATE OR REPLACE FUNCTION update_follow()
    RETURNS TRIGGER AS $$
BEGIN
    -- Update the followed_products for the user
    INSERT INTO follow (user_id, user_name, followed_products)
    VALUES (
               NEW.user_id,
               (SELECT name FROM users WHERE id = NEW.user_id),
               (
                   SELECT jsonb_agg(
                                  jsonb_build_object(
                                          'product_id', product_id,
                                          'product_name', product_name,
                                          'product_description', product_description,
                                          'product_price', product_price
                                  )
                          )
                   FROM follow_list
                   WHERE user_id = NEW.user_id
               )
           )
    ON CONFLICT (user_id) DO UPDATE
        SET followed_products = EXCLUDED.followed_products;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER after_follow_list_change
    AFTER INSERT OR UPDATE OR DELETE ON follow_list
    FOR EACH ROW
EXECUTE FUNCTION update_follow();

--EXAMPLE OF INSERTING
insert into follow_list(user_id, product_id)
values(1, 4);
