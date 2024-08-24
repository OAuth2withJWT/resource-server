CREATE TABLE cards (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL,  
    card_number VARCHAR(255) NOT NULL,
    current_balance DECIMAL(10, 2) NOT NULL
);