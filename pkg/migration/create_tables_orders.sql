CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INT NOT NULL,
                                      total_amount NUMERIC(10, 2) NOT NULL,
                                      delivery_addr TEXT NOT NULL,
                                      status VARCHAR(255) NOT NULL,
                                      created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

drop table orders;

CREATE TABLE IF NOT EXISTS order_products (
                                              id SERIAL PRIMARY KEY,
                                              order_id INT NOT NULL,
                                              product_id INT NOT NULL,
                                              quantity INT NOT NULL,
                                              FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                                              FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);