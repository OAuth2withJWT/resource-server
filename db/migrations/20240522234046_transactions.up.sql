CREATE TYPE expense_category AS ENUM (
    'monthly',
    'groceries',
    'healthcare',
    'clothing',
    'entertainment',
    'dining',
    'transport',
    'utilities',
    'transfer'
);

CREATE TYPE transaction_type AS ENUM (
    'income',
    'expense'
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,                
    time TIMESTAMP NOT NULL,                  
    amount DECIMAL(10, 2) NOT NULL,      
    expense_category expense_category,
    transaction_type transaction_type NOT NULL,
    location VARCHAR(100)
);
