CREATE TABLE recommended_products (
    product_id BIGINT NOT NULL,
    recommended_id BIGINT NOT NULL,
    PRIMARY KEY (product_id, recommended_id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_recommended FOREIGN KEY (recommended_id) REFERENCES products(id) ON DELETE CASCADE
);

ALTER TABLE recommended_products
ADD COLUMN created_at TIMESTAMP DEFAULT NOW();


CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    size_id INT REFERENCES sizes(id),
    variant_id BIGINT REFERENCES variants(id),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- index supaya pencarian cepat
CREATE INDEX idx_carts_user ON carts(user_id);
CREATE INDEX idx_carts_product ON carts(product_id);


CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    fullname VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    address VARCHAR(250),
    payment_method_id INT NOT NULL,
    shipping_id INT NOT NULL,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    total NUMERIC(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id),
    FOREIGN KEY (shipping_id) REFERENCES shippings(id)
);

CREATE TABLE transaction_items (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    variant_id INT,      
    size_id INT,        
    quantity INT NOT NULL,
    subtotal NUMERIC(10,2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    FOREIGN KEY (size_id) REFERENCES sizes(id)
);



