drop table if exists users;
drop table if exists products;
drop table if exists shop;
drop table if exists catalog_and_products;

INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'shop:read') FROM users;

update users set activated=True where email='az_tau@kbtu.kz';

INSERT INTO users_permissions
VALUES (
           (SELECT id FROM users WHERE email = 'name_4@kbtu.kz'),
           (SELECT id FROM permissions WHERE code = 'shop:write')
       );
