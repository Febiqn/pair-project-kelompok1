create database rental;
use rental;

CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    membership_status ENUM('ACTIVE','INACTIVE') DEFAULT 'INACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE playstations (
    ps_id INT AUTO_INCREMENT PRIMARY KEY,
    ps_name VARCHAR(50) NOT NULL,
    condition_status ENUM('AVAILABLE', 'RENTED', 'BROKEN') DEFAULT 'AVAILABLE'
);

CREATE TABLE rentals (
    rental_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    ps_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    duration_hours INT NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('ONGOING', 'COMPLETED') DEFAULT 'ONGOING',

    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (ps_id) REFERENCES playstations(ps_id)
);

CREATE TABLE billing (
    bill_id INT AUTO_INCREMENT PRIMARY KEY,
    rental_id INT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (rental_id) REFERENCES rentals(rental_id)
);


