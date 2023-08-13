create table if not exists users (
    user_id varchar(36) PRIMARY KEY,
    username varchar(255),
    email varchar(255) unique,
    password varchar(255),
    role ENUM('admin', 'user'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );

create table if not exists product (
    product_id varchar(36) PRIMARY KEY,
    category_id varchar(36),
    name varchar(255),
    description varchar(255),
    price int,
    stock int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );
create table if not exists product_categories (
    category_id varchar(36) PRIMARY KEY,
    name varchar(255),
    description varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );

create table if not exists orders (
    order_id varchar(36) PRIMARY KEY,
    user_id varchar(36),
    total_amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );
create table if not exists order_items (
    order_item_id varchar(36) PRIMARY KEY,
    order_id varchar(36),
    product_id varchar(36),
    quantity int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );
create table if not exists carts (
    cart_id varchar(36) PRIMARY KEY,
    user_id varchar(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );
create table if not exists cart_items (
    cart_item_id varchar(36) PRIMARY KEY,
    cart_id varchar(36),
    product_id varchar(36),
    quantity int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(36)
    );

CREATE INDEX idx_product_category_id ON product (category_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_order_items_order_product ON order_items (order_id, product_id);