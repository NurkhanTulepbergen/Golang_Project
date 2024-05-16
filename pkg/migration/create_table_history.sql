create table if not exists history
(
    user_id int,
    user_name varchar(255),
    orders_list jsonb,
    foreign key(user_id) references users(id)
);

INSERT INTO history (user_id, user_name, orders_list)
SELECT o.user_id, h.user_name, json_agg(json_build_object('id', o.id, 'products', o.products, 'total_amount', o.total_amount, 'delivery_addr', o.delivery_addr, 'status', o.status, 'created_at', o.created_at))
FROM orders o
         INNER JOIN history h ON o.user_id = h.user_id
WHERE o.user_id = 1
GROUP BY o.user_id, h.user_name;

