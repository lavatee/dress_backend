CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('customer', 'buyer', 'admin'))
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    main_photo_url VARCHAR(255),
    category_id INT NOT NULL,
    collection VARCHAR(255),
    color VARCHAR(50),
    price INT,
    description VARCHAR(1023)
);

CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('photo', 'video')),
    url VARCHAR(255) NOT NULL UNIQUE,
    product_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS sizes (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    name VARCHAR(15) NOT NULL,
    amount BIGINT,
    UNIQUE(product_id, name)
);

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    user_id INT NOT NULL,
    rating INT NOT NULL,
    comment VARCHAR(1000),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, user_id)
);

CREATE TABLE IF NOT EXISTS liked_products (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    user_id INT NOT NULL,
    UNIQUE(product_id, user_id)
);

CREATE TABLE IF NOT EXISTS products_in_cart (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    user_id INT NOT NULL,
    size VARCHAR(15) NOT NULL,
    amount INT NOT NULL,
    existence BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (user_id, size)
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('delivery', 'pickup')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'created', 'delivered_to_shop', 'issued', 'sent_to_customer', 'delivered_to_customer')),
    shop_point VARCHAR(255),
    user_id INT NOT NULL,
    payment_id VARCHAR(255),
    order_price INT,
    delivery_id VARCHAR(255),
    pickup_id VARCHAR(255),
    delivery_price INT,
    delivery_address VARCHAR(255),
    delivery_index INT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ordered_products (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL, 
    size VARCHAR(15) NOT NULL,
    amount BIGINT NOT NULL,
    price INT NOT NULL,
    product_name VARCHAR(255) NOT NULL
);

ALTER TABLE reviews ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE reviews ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE media ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE sizes ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE liked_products ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE products_in_cart ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories(id);
ALTER TABLE liked_products ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE products_in_cart ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE orders ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE ordered_products ADD FOREIGN KEY (order_id) REFERENCES orders(id);
