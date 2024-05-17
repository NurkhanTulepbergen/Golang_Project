create table if not exists history
(
    user_id INT primary key,
    user_name VARCHAR(255),
    orders_list JSONB
);

CREATE OR REPLACE FUNCTION update_history()
    RETURNS TRIGGER AS $$
BEGIN
    -- Update the followed_products for the user
    INSERT INTO history (user_id, user_name, orders_list)
    VALUES (
               NEW.user_id,
               (SELECT name FROM users WHERE id = NEW.user_id),
               (
                   SELECT jsonb_agg(
                                  jsonb_build_object(
                                          'order_id', id,
                                          'user_id', user_id,
                                          'total_amount', total_amount,
                                          'delivery_addr', delivery_addr,
                                          'status', status,
                                          'created_at', created_at
                                  )
                          )
                   FROM orders
                   WHERE user_id = NEW.user_id
               )
           )
    ON CONFLICT (user_id) DO UPDATE
        SET orders_list = EXCLUDED.orders_list;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER after_add_order
    AFTER INSERT OR UPDATE OR DELETE ON orders
    FOR EACH ROW
EXECUTE FUNCTION update_history();

-- INSERT INTO history (user_id, user_name, orders_list)
-- SELECT o.user_id, h.user_name, json_agg(json_build_object('id', o.id, 'products', o.products, 'total_amount', o.total_amount, 'delivery_addr', o.delivery_addr, 'status', o.status, 'created_at', o.created_at))
-- FROM orders o
--          INNER JOIN history h ON o.user_id = h.user_id
-- WHERE o.user_id = 1
-- GROUP BY o.user_id, h.user_name;

drop table history;