CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,                
    time TIMESTAMP NOT NULL,                  
    amount DECIMAL(10, 2) NOT NULL,      
    category VARCHAR(100),
    location VARCHAR(100)
);