# Pet Agency API Documentation

## Overview
The Pet Agency API allows users to register, login, add pet listings, edit their listings, and toggle the status of their listings. The API is secured with JWT (JSON Web Token) and users can only edit the data they created. There are two public endpoints that return all available pets for adoption and the details of a single animal.

## Technology Stack
- **Language:** Golang
- **Framework:** Echo
- **Database:** PostgreSQL

## API Endpoints

### Public Endpoints

#### 1. Get All Available Pets
- **Endpoint:** `GET /api/pets`
- **Description:** Returns a list of all available pets for adoption.
- **Response:**
  ```json
  [
      {
          "id": 1,
          "name": "Bella",
          "species": "Dog",
          "breed": "Labrador",
          "age": 3,
          "status": "Available"
      },
      ...
  ]
  ```

#### 2. Get Pet Details
- **Endpoint:** `GET /api/pets/:id`
- **Description:** Returns details of a single pet.
- **Response:**
  ```json
  {
      "id": 1,
      "name": "Bella",
      "species": "Dog",
      "breed": "Labrador",
      "age": 3,
      "description": "Friendly and energetic",
      "status": "Available"
  }
  ```

### Protected Endpoints

#### 3. Register
- **Endpoint:** `POST /api/auth/register`
- **Description:** Registers a new user.
- **Request:**
  ```json
  {
      "username": "user1",
      "email": "user1@example.com",
      "password": "password123"
  }
  ```
- **Response:**
  ```json
  {
      "message": "User registered successfully"
  }
  ```

#### 4. Login
- **Endpoint:** `POST /api/auth/login`
- **Description:** Logs in a user and returns a JWT token.
- **Request:**
  ```json
  {
      "email": "user1@example.com",
      "password": "password123"
  }
  ```
- **Response:**
  ```json
  {
      "token": "JWT_TOKEN"
  }
  ```

#### 5. Add Pet Listing
- **Endpoint:** `POST /api/pets`
- **Description:** Adds a new pet listing. Requires JWT.
- **Request:**
  ```json
  {
      "name": "Bella",
      "species": "Dog",
      "breed": "Labrador",
      "age": 3,
      "description": "Friendly and energetic"
  }
  ```
- **Response:**
  ```json
  {
      "message": "Pet listing created successfully"
  }
  ```

#### 6. Edit Pet Listing
- **Endpoint:** `PUT /api/pets/:id`
- **Description:** Edits an existing pet listing. Requires JWT and ownership.
- **Request:**
  ```json
  {
      "name": "Bella",
      "species": "Dog",
      "breed": "Labrador",
      "age": 3,
      "description": "Friendly and very energetic"
  }
  ```
- **Response:**
  ```json
  {
      "message": "Pet listing updated successfully"
  }
  ```

#### 7. Toggle Pet Status
- **Endpoint:** `PATCH /api/pets/:id/status`
- **Description:** Toggles the status of a pet listing (Available/Adopted). Requires JWT and ownership.
- **Request:**
  ```json
  {
      "status": "Adopted"
  }
  ```
- **Response:**
  ```json
  {
      "message": "Pet status updated successfully"
  }
  ```

## Database Schema

### Users Table

| Column     | Type         | Constraints               |
|------------|--------------|---------------------------|
| id         | SERIAL       | PRIMARY KEY               |
| username   | VARCHAR(255) | NOT NULL, UNIQUE          |
| email      | VARCHAR(255) | NOT NULL, UNIQUE          |
| password   | VARCHAR(255) | NOT NULL                  |

### Pets Table

| Column       | Type         | Constraints               |
|--------------|--------------|---------------------------|
| id           | SERIAL       | PRIMARY KEY               |
| name         | VARCHAR(255) | NOT NULL                  |
| species      | VARCHAR(255) | NOT NULL                  |
| breed        | VARCHAR(255) |                           |
| age          | INT          |                           |
| description  | TEXT         |                           |
| status       | VARCHAR(50)  | NOT NULL                  |
| user_id      | INT          | FOREIGN KEY (users.id)    |

## Security

### JWT Authentication
- JWT is used to secure the API endpoints.
- Users must include the JWT token in the Authorization header for protected endpoints.

## Implementation Details

### Middleware
- **JWT Middleware:** Ensures that the token is valid and the user is authenticated.
- **Ownership Middleware:** Ensures that the user can only edit or toggle the status of pets they created.

### Error Handling
- Proper error messages and status codes are returned for invalid requests and unauthorized actions.

---

Feel free to adjust any details or add any additional sections as needed.