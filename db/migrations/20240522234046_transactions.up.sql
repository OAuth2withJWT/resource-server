CREATE TYPE transaction_category AS ENUM (
    'monthly',
    'groceries',
    'healthcare',
    'rent',
    'utilities',
    'savings',
    'transportation',
    'clothing',
    'personal_care'
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,                
    time TIMESTAMP NOT NULL,                  
    amount DECIMAL(10, 2) NOT NULL,      
    category transaction_category NOT NULL,
    location VARCHAR(100) NOT NULL
);
