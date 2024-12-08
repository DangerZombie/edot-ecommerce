CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    password TEXT
);

-- CREATE TABLE IF NOT EXISTS products (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT,
--     price REAL,
--     stock INTEGER
-- );

-- CREATE TABLE IF NOT EXISTS orders (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     user_id INTEGER,
--     product_id INTEGER,
--     quantity INTEGER,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE IF NOT EXISTS shops (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT,
--     owner_id INTEGER
-- );

-- CREATE TABLE IF NOT EXISTS warehouses (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     location TEXT,
--     capacity INTEGER
-- );