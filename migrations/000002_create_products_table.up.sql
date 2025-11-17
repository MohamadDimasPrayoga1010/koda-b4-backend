CREATE TABLE variants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    additional_price INT DEFAULT 0
);


CREATE TABLE categories(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
); 

CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    description VARCHAR(250),
    stock INT,
    category_id BIGINT REFERENCES categories(id),
    base_price NUMERIC,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT now()
); 

ALTER TABLE products ALTER COLUMN deleted_at SET DEFAULT NULL;
ALTER TABLE products
ADD COLUMN is_favorite BOOLEAN DEFAULT false;

CREATE TABLE product_variants (
    id SERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id BIGINT NOT NULL REFERENCES variants(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE products_categories(
    product_id BIGINT REFERENCES products(id),
    category_id BIGINT REFERENCES categories(id)
);

CREATE TABLE product_images(
    id SERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id),
    image TEXT,
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT now()
);
 
CREATE TABLE sizes(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    additional_price NUMERIC
);

CREATE TABLE product_sizes(
    product_id BIGINT REFERENCES products(id),
    size_id INT REFERENCES sizes(id)
);