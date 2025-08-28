-- Initial schema for Dunhayat Coffee Roastery
-- Migration: 20241201000000_initial_schema.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255),
    verified SMALLINT DEFAULT 0 CHECK (verified >= 0 AND verified <= 3),
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Addresses table
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE products (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    description TEXT,
    price INTEGER NOT NULL CHECK (price > 0),
    image_url TEXT,
    category SMALLINT NOT NULL,
    in_stock INTEGER DEFAULT 0 CHECK (in_stock >= 0),
    weight DECIMAL(8,2),
    origin VARCHAR(100),
    roast_level VARCHAR(50),
    bitterness INTEGER CHECK (bitterness >= 1 AND bitterness <= 5),
    body INTEGER CHECK (body >= 1 AND body <= 5),
    acidity INTEGER CHECK (acidity >= 1 AND acidity <= 5),
    sweetness INTEGER CHECK (sweetness >= 1 AND sweetness <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sales (orders) table
CREATE TABLE sales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    tracking_code VARCHAR(100),
    total_price INTEGER NOT NULL CHECK (total_price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sale items table
CREATE TABLE sale_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sale_id UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cart reservations table (for 10-minute product reservations)
CREATE TABLE cart_reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table for user authentication
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_addresses_user_id ON addresses(user_id);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_in_stock ON products(in_stock);
CREATE INDEX idx_sales_user_id ON sales(user_id);
CREATE INDEX idx_sales_status ON sales(status);
CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_cart_reservations_user_id ON cart_reservations(user_id);
CREATE INDEX idx_cart_reservations_expires_at ON cart_reservations(expires_at);
CREATE INDEX idx_sessions_token ON sessions(token);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON sales
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
