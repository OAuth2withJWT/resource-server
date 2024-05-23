CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,                
    date DATE NOT NULL,                  
    amount DECIMAL(10, 2) NOT NULL,      
    category VARCHAR(100)
);