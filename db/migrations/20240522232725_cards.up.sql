CREATE TABLE cards (
    card_id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL,  
    card_number VARCHAR(20) NOT NULL,
    current_balance DECIMAL(10, 2) NOT NULL
);