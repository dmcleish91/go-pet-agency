CREATE DATABASE pet-agency;

CREATE TABLE pet-agency.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashedpassword VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pet-agency.pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    species VARCHAR(255) NOT NULL,
    breed VARCHAR(255),
    age INT,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    gender VARCHAR(50),
    size VARCHAR(50),
    color VARCHAR(255),
    weight DECIMAL(5,2),
    vaccination_status BOOLEAN,
    spayed BOOLEAN,
    microchipped BOOLEAN,
    rescue_story TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES pet-agency.users(id)
);

-- Create a New User
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3);  -- Assuming $1, $2, and $3 are placeholders for username, email, and password

-- Login Details
SELECT id, username, password
FROM users
WHERE email = $1;  -- Assuming $1 is the placeholder for the user's email

-- Update User Email
UPDATE users
SET email = $1  -- Assuming $1 is the placeholder for the new email address
WHERE id = $2;  -- Assuming $2 is the placeholder for the user's ID

-- Check if a email address exists in the user table
SELECT COUNT(*)
FROM users
WHERE email = 'email@example.com';

-- Update Username
UPDATE users
SET username = $1  -- Assuming $1 is the placeholder for the new username
WHERE id = $2;  -- Assuming $2 is the placeholder for the user's ID

-- Get All Available Pets
SELECT id, name, species, breed, age, status, gender, size, color, weight, vaccination_status, spayed, microchipped, rescue_story, created_at, user_id
FROM pets
WHERE status = 'Available';

-- Get Pet Details
SELECT id, name, species, breed, age, description, status, gender, size, color, weight, vaccination_status, spayed, microchipped, rescue_story, created_at
FROM pets
WHERE id = $1;  -- Assuming $1 is the placeholder for the pet's ID

-- Add Pet Listing
INSERT INTO pets (name, species, breed, age, description, status, gender, size, color, weight, vaccination_status, spayed, microchipped, rescue_story, user_id)
VALUES ($1, $2, $3, $4, $5, 'Available', $6, $7, $8, $9, $10, $11, $12, $13, $14);  -- Assuming $1 to $14 are placeholders for pet details and user_id

-- Edit Pet Listing
UPDATE pets
SET name = $1, species = $2, breed = $3, age = $4, description = $5, gender = $6, size = $7, color = $8, weight = $9, vaccination_status = $10, spayed = $11, microchipped = $12, rescue_story = $13
WHERE id = $14 AND user_id = $15;  -- Assuming $1 to $15 are placeholders for pet details, pet ID, and user ID

-- Toggle Pet Status
UPDATE pets
SET status = $1  -- Assuming $1 is the placeholder for the new status
WHERE id = $2 AND user_id = $3;  -- Assuming $2 is the pet ID and $3 is the user ID


