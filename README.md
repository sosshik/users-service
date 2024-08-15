# Users Management Microservice

## Build Status

![Build Status](https://img.shields.io/github/actions/workflow/status/sosshik/users-service/go.yml?branch=main)


## Overview

This microservice manages access to users, allowing you to add, modify, remove, and retrieve a paginated list of users with filtering capabilities. The service is built using Go and employs in-memory storage to manage user data.

## Features

- **Add a new User:** Add a new user with required attributes.
- **Modify an existing User:** Update existing user details using their ID.
- **Remove a User:** Delete a user using their ID.
- **Retrieve Users:** Fetch a paginated list of users, with optional filtering by specific criteria (e.g., country).
- **Health Check:** A simple health check endpoint to monitor service status.

## API Documentation 

API Documentation is available at http://localhost:8090/swagger/index.html

## Tests
Unit tests are provided for the core functionality of the service. You can run the tests using:
```bash
go test -v ./...
```

## Setup Instructions
### Prerequisites
- Go 1.20+
- Docker (if you prefer running the application in a container)

### Running the Application
**Option 1: Running Locally**
1. Clone the Repository:

```bash
git clone https://github.com/sosshik/users-service
cd users-service
```

2. Install Dependencies:
```bash
go mod download
```
3. Build and Run the Application:

```bash
go build -o users-service cmd/users-service/main.go
./users-service
```
4. Access the API:

The service will be running on http://localhost:8090.

**Option 2: Running with Docker**
 1. Build the Docker Image:

```bash
docker build -t users-service .
```

2. Run the Docker Container:
```bash
docker run -p 8090:8090 users-service
```
3. Access the API:

The service will be running on http://localhost:8090.

## Logging
The service uses `logrus` for structured logging. Logs provide insights into the service's operations, including warnings and errors.

## Documentation
Swagger documentation is provided and can be accessed at http://localhost:8090/swagger/index.html once the service is running. For generating docs I used `swaggo`

## Assumptions and Choices
- **In-Memory Storage:** For simplicity, the service uses in-memory storage. This decision was made to align with the requirement to not use a SQL database and to focus on the core functionality.
- **Go with Echo Framework:** Echo was chosen for its simplicity and performance in building HTTP APIs. It also provides easy integration with middleware and is a common choice in Go-based microservices.
- **Logging:** `logrus` was chosen for structured logging to easily integrate with centralized logging systems in a production environment.
- **Swagger for API Documentation:** Swagger provides a simple way to document and test the API, making it easier for developers and testers to interact with the service.
- **UUID:** UUIDs are used for unique user identification to ensure each user has a globally unique identifier.
- **`jinzhu/copier`:** This package is used to cast models and DTOs efficiently, simplifying the process of copying data between structures.
- **`go-ozzo/ozzo-validation`:** This package is used for validating requests, providing a robust way to ensure incoming data meets specified requirements.

## Possible Extensions and Improvements
- **Persistent Storage:** Replace in-memory storage with a persistent database like PostgreSQL or MongoDB for long-term data storage.
- **Authentication:** Implement authentication mechanisms like OAuth2 or JWT for securing endpoints.
- **Improved Pagination and Filtering:** Extend filtering options to include more fields, and refine pagination to handle large datasets efficiently.
- **Load Balancing & Scaling:** Integrate with a load balancer and run multiple instances of the service for scalability in a production environment.
- **Configuration Management:** Introduce a configuration management system to manage application settings and environment variables, allowing for easier updates and better environment separation (e.g., using configuration files or environment variables).
- **Caching:** Implement caching to reduce the load on the service and improve response times. Caching frequently accessed data (e.g., user details) can significantly enhance performance, especially for read-heavy operations.