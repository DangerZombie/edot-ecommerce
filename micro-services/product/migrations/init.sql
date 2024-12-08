CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL,
    stock INTEGER NOT NULL
);

-- Insert dummy data into products table if not exists
INSERT INTO products (name, description, price, stock) 
SELECT 'Product A', 'Description of Product A', 100.0, 50
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Product A');

INSERT INTO products (name, description, price, stock) 
SELECT 'Product B', 'Description of Product B', 200.0, 20
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Product B');

INSERT INTO products (name, description, price, stock) 
SELECT 'Product C', 'Description of Product C', 150.0, 30
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'Product C');