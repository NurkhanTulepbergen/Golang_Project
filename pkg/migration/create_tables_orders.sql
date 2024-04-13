CREATE TABLE Orders (
                       ID VARCHAR(255) PRIMARY KEY,
                       "user" VARCHAR(255),
                       TotalAmount NUMERIC(10, 2),
                       DeliveryAddr VARCHAR(255),
                       Status VARCHAR(50),
                       CreatedAt TIMESTAMP
);
