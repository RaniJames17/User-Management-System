
# User Management System

**DIRECTORY STRUCTURE**
user-management-system/
├── controllers/
│   └── user.go                   # Contains logic for sign-up, sign-in, reset password, etc.
├── database/
│   └── connection.go             # Database connection setup
│   └── query.go                  # responsible for executing SQL queries in the MySQL database
├── models/
│   └── user.go                   # Defines the User model
│   └── passwordreset.go          # Defines the PasswordReset model
├── routes/
│   └── routes.go                 # Defines all the routes for your API
├── middlewares/
│   └── auth.go                   # Middleware for token authentication
├── utils/
│   └── hash.go                   # Helper functions - password hashing 
│   └── token.go                  # Helper functions - token generation
├── Dockerfile                    # Docker setup for the application
├── docker-compose.yml            # Docker Compose configuration for MySQL and the app
├── go.mod                        # Go module file
├── go.sum                        # Go checksum file
└── main.go                       # Entry point for the application


## Project Description
The **User Management System** is a web application that allows users to register, sign in, and manage their passwords. The system supports token-based authentication, password reset functionality, and input validation. It provides an easy-to-use API for user management, including functionalities like user registration, login, and password recovery.

### Features:
- **User Sign-Up**: Register new users with email, password, and name. The system validates the email format and password strength.
- **User Sign-In**: Authenticate users via email and password, and issue a JSON Web Token (JWT) for session management.
- **Forgot Password**: Users can request a password reset token that is sent to their registered email.
- **Reset Password**: Users can reset their passwords using the reset token sent to their email.

## Installation Instructions

### Prerequisites
- Docker
- MySQL (Dockerized or installed locally)
- Go (for building the app)

### Setting Up Docker and MySQL

#### 1. Docker Setup
To get the application running with Docker, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repository-url.git
   cd your-repository-name
   ```

2. Create a Docker network for your services:
   ```bash
   docker network create user-management-network
   ```

3. Build the Docker containers:
   - In the root directory of your project, create a `docker-compose.yml` file:
     ```yaml
     version: '3.7'

     services:
       db:
         image: mysql:8.0
         container_name: mysql-container
         environment:
           MYSQL_ROOT_PASSWORD: rootpassword
           MYSQL_DATABASE: user_management
           MYSQL_USER: user
           MYSQL_PASSWORD: userpassword
         ports:
           - "3306:3306"
         networks:
           - user-management-network

       app:
         build: .
         container_name: user-management-app
         environment:
           DB_HOST: db
           DB_USER: user
           DB_PASSWORD: userpassword
           DB_NAME: user_management
         ports:
           - "8080:8080"
         depends_on:
           - db
         networks:
           - user-management-network

     networks:
       user-management-network:
         external: true
     ```

4. Start the Docker containers:
   ```bash
   docker-compose up --build
   ```

#### 2. MySQL Setup
- The Docker Compose file automatically sets up a MySQL container. The application connects to this MySQL database to store user data.
- The database has the following schema:
  ```sql
  CREATE DATABASE user_management;
  
  CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
  );
  
  CREATE TABLE password_resets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
  ```

### Running the Application Locally

1. **Set up environment variables**:
   - Create a `.env` file in the root directory with the following variables:
     ```env
     DB_HOST=localhost
     DB_USER=user
     DB_PASSWORD=userpassword
     DB_NAME=user_management
     ```

2. **Install dependencies**:
   Make sure you have Go installed and then run:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. The application will be available at `http://localhost:8080`.

---

## API Documentation

### 1. Health Check
- **Endpoint**: `/api/health`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "message": "API is up and running"
  }
  ```

### 2. User Sign-Up
- **Endpoint**: `/api/signup`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "strongpassword123"
  }
  ```
- **Response**:
  ```json
  {
    "message": "User created successfully"
  }
  ```

### 3. User Sign-In
- **Endpoint**: `/api/signin`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "strongpassword123"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Sign-In successful",
    "token": "JWT_TOKEN",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com"
    }
  }
  ```

### 4. Forgot Password
- **Endpoint**: `/api/forgot-password`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com"
  }
  ```
- **Response**:
  ```json
  {
    "message": "If the email is valid, you will receive a reset link."
  }
  ```

### 5. Reset Password
- **Endpoint**: `/api/reset-password`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "reset_token": "RESET_TOKEN",
    "new_password": "newpassword123"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Password reset successfully"
  }
  ```

---

## Team Members & Contributions

### 1. **Rani James**
- **Contributions**:
  - Implemented forgot password, and reset password features.
  - Set up the MySQL database schema and integrated it with the application.
  - Setup github

### 2. **Harninder Singh**
- **Contributions**:
  - Implemented user sign-up, sign-in
  - Added required tables 
  - Implemented token

### 3. **Jatin Sachdev**
- **Contributions**:
  - Set up Docker and Docker Compose configurations for easy deployment of the application.
  - Configured the MySQL database for Docker container integration.

### 3. **Ishaben Bhatt**
- **Contributions**:
  - Wrote API documentation and ensured input validation.
  - Created Readme.md, tested the application
