CREATE TABLE if not exists users
(
    user_id  serial PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email, password)
VALUES ('a', 'a@gmail.com', 'aaa');

CREATE TABLE if not exists shipments (
    shipment_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    modelType VARCHAR(50) NOT NULL,
    projectName VARCHAR(255) NOT NULL,
    algorithm VARCHAR(255) NOT NULL,
    targetColumn VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE if not exists downloaded_files (
    file_id SERIAL PRIMARY KEY,
    shipment_id INT NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shipment_id) REFERENCES shipments(shipment_id)
);

CREATE TABLE if not exists model_files (
    file_id SERIAL PRIMARY KEY,
    shipment_id INT NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shipment_id) REFERENCES shipments(shipment_id)
);

CREATE TABLE if not exists model_metrics (
    metric_id SERIAL PRIMARY KEY,
    file_id INT NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    metric_value FLOAT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES model_files(file_id)
);