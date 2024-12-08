CREATE TABLE IF NOT EXISTS shops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT
);

INSERT INTO shops (name, description)
SELECT 'Shop A', 'Description of Shop A'
WHERE NOT EXISTS (SELECT 1 FROM shops WHERE name = 'Shop A');
