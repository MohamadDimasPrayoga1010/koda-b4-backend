CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE, 
    password VARCHAR(100) NOT NULL,
    role VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

ALTER TABLE users
ADD COLUMN reset_token TEXT,
ADD COLUMN reset_expires TIMESTAMP;
ALTER TABLE users
ADD COLUMN reset_otp VARCHAR(25);

CREATE TABLE profile(
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id),
    image VARCHAR(250),
    phone VARCHAR(50),
    address VARCHAR(250),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE forgot_password(
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id),
    token VARCHAR(100),
    expires_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
); 