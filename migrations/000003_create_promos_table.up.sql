CREATE TABLE promos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    discount FLOAT NOT NULL,
    start TIMESTAMP NOT NULL,
    "end" TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE TABLE product_promos (
    id SERIAL PRIMARY KEY,
    promo_id BIGINT REFERENCES promos(id),
    product_id BIGINT REFERENCES products(id)
); 
