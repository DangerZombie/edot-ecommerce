CREATE TABLE IF NOT EXISTS warehouses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    status TEXT
);

INSERT INTO warehouses (name, status)
SELECT 'Warehouse A', 'active'
WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Warehouse A');

INSERT INTO warehouses (name, status)
SELECT 'Warehouse B', 'active'
WHERE NOT EXISTS (SELECT 1 FROM warehouses WHERE name = 'Warehouse B');

CREATE TABLE IF NOT EXISTS stocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    warehouse_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    UNIQUE (warehouse_id, product_id)  -- Add UNIQUE constraint on warehouse_id and product_id
);

-- Insert initial stock data based on product stock

-- For Warehouse A
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 25  -- Assign specific stock quantity for Product A
FROM warehouses w
JOIN products p ON p.name = 'Product A'
WHERE w.name = 'Warehouse A'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- For Warehouse B
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 25  -- Assign specific stock quantity for Product A
FROM warehouses w
JOIN products p ON p.name = 'Product A'
WHERE w.name = 'Warehouse B'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- Repeat for Product B and Product C
-- For Warehouse A
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 5  -- Assign specific stock quantity for Product B
FROM warehouses w
JOIN products p ON p.name = 'Product B'
WHERE w.name = 'Warehouse A'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- For Warehouse B
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 15  -- Assign specific stock quantity for Product B
FROM warehouses w
JOIN products p ON p.name = 'Product B'
WHERE w.name = 'Warehouse B'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- For Warehouse A
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 10  -- Assign specific stock quantity for Product C
FROM warehouses w
JOIN products p ON p.name = 'Product C'
WHERE w.name = 'Warehouse A'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- For Warehouse B
INSERT INTO stocks (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 20  -- Assign specific stock quantity for Product C
FROM warehouses w
JOIN products p ON p.name = 'Product C'
WHERE w.name = 'Warehouse B'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;