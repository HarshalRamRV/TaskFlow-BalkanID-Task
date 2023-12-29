# TaskFlow (Task Management System)

TaskFlow is a containerized task management system inspired by Atlassian Jira. It provides user authentication, authorization, and access management. Users can register, authenticate, manage their accounts, and interact with tasks. The system is built using React.js, MySQL, and Docker for containerization.


## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Frontend](#frontend)
- [Docker](#docker)
- [Security](#security)

## Getting Started

### Prerequisites

Before you begin, ensure you have the following tools installed:

- Node.js and npm: [Download Node.js](https://nodejs.org/)
- MySQL: [Download MySQL](https://dev.mysql.com/downloads/)
- Docker: [Download Docker](https://www.docker.com/get-started)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/HarshalRamRV/TaskFlow-BalkanID-Task
   cd TaskFlow-BalkanID-Task
   ```

2. Install backend and frontend dependencies:

   ```bash
   # Install backend dependencies (Go)
   cd backend
   go get ./...

   # Install frontend dependencies
   cd ../frontend
   npm install
   ```


3. Set up the database:

    ```sql
    CREATE DATABASE task_management;

    USE task_management;

    CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL,
        roleId VARCHAR(255) NOT NULL,
        groupId VARCHAR(255) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        deactivated TINYINT DEFAULT 0,
        roleId VARCHAR(255) NOT NULL,
        groupId VARCHAR(255) NOT NULL
    );
    ```

4. Initialize the database schema and seed data:

   ```bash
   cd backend
   go run setup.go
   ```

5. Build the frontend:

   ```bash
   cd frontend
   npm run build
   ```

6. Start the application using Docker Compose:

   ```bash
   docker-compose up
   ```

Your application should be accessible at `http://localhost`.


## Frontend

The frontend is built with React.js and includes user-friendly interfaces for task management. It uses Tailwind CSS for styling.

## Docker

The application is containerized using Docker Compose, making it easy to deploy. The Docker configuration is in the `docker-compose.yml` file.

## Security

- User authentication and authorization are implemented.
- Protection against common vulnerabilities like SQL injection is provided.
