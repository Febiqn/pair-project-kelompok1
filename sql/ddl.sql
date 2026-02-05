CREATE DATABASE IF NOT EXISTS rental;
USE rental;

CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    membership_status ENUM('ACTIVE','INACTIVE') DEFAULT 'INACTIVE',
    membership_number VARCHAR(20) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS playstations (
    ps_id INT AUTO_INCREMENT PRIMARY KEY,
    ps_name VARCHAR(50) NOT NULL,
    condition_status ENUM('AVAILABLE','RENTED','BROKEN') DEFAULT 'AVAILABLE'
);

CREATE TABLE IF NOT EXISTS rentals (
    rental_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    user_name VARCHAR(100) NOT NULL,
    ps_id INT NOT NULL,
    ps_name VARCHAR(50) NOT NULL,
    membership_status ENUM('ACTIVE','INACTIVE') NOT NULL,
    start_time DATETIME NOT NULL,
    duration_hours INT NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('ONGOING','COMPLETED') DEFAULT 'ONGOING',
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (ps_id) REFERENCES playstations(ps_id)
);

CREATE TABLE IF NOT EXISTS billing (
    bill_id INT AUTO_INCREMENT PRIMARY KEY,
    rental_id INT NOT NULL,
    total_amount DECIMAL(10,2),
    bill_status ENUM('UNPAID','PAID') DEFAULT 'UNPAID',
    paid_at DATETIME NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rental_id) REFERENCES rentals(rental_id)
);
