# Workout Tracker API

## Overview
Workout Tracker API is a RESTful API designed to help users create, manage, and track their workout plans. It provides endpoints for creating, updating, deleting, and listing workout plans, with secure user authentication using JSON Web Tokens (JWT).

## Features
- **User Authentication**: Secure user registration and login with JWT-based authentication.
- **Workout Plan Management**: Create, update, delete, and list workout plans.
- **RESTful Design**: Simple and intuitive endpoints for seamless integration.

## Endpoints

### Workout Plans
- **POST /workout-plans/**  
  Create a new workout plan.  
  **Requires Authentication**: Yes  
  **Request Body**: JSON with workout plan details (e.g., ExerciseName, Repetition, Sets, Weights).  
  **Response**: JSON error message if error is encountered

- **PUT /workout-plans/**  
  Update an existing workout plan.  
  **Requires Authentication**: Yes  
  **Request Body**: JSON with updated workout plan details.  
  **Response**: JSON error message if error is encountered

- **DELETE /workout-plans/**  
  Delete a workout plan.  
  **Requires Authentication**: Yes  
  **Request Body**: JSON with workout plan name.  
  **Response**: JSON error message if error is encountered

- **GET /workouts**  
  List all workout plans for the authenticated user.  
  **Requires Authentication**: Yes  
  **Response**: JSON array of workout plans or error message.

### Authentication
- **POST /auth/register**  
  Register a new user.  
  **Request Body**: JSON with user details (e.g., username, email, password).  
  **Response**: JSON with user details and JWT token or error message.

- **POST /auth/login**  
  Log in an existing user.  
  **Request Body**: JSON with credentials (e.g., email, password).  
  **Response**: JSON with JWT token or error message.

## Authentication
The API uses **JWT (JSON Web Tokens)** for secure authentication. Include the JWT in the `Authorization` header for protected endpoints:  
```
Authorization: Bearer <your-jwt-token>
```

## Getting Started

### Prerequisites
- **Go**: Ensure Go is installed (version 1.16 or later recommended).
- **Database**: The API requires a configured database (e.g., PostgreSQL, MySQL). Ensure the database server is running and accessible.
- **Environment Variables**: Set up a `.env` file with the necessary environment variables (see Configuration section below).

### Configuration
1. **Database Setup**:
   - Create a database for the API (e.g., `workout_tracker_db`).
   - Configure the database connection by setting the following environment variables in a `.env` file in the project root:
     ```
     DB_USER=<your-database-username>
     DB_NAME=<your-database-name>
     DB_PASSWORD=<your-database-password>
     ```
   - Ensure the database user has appropriate permissions to create and modify tables.

2. **JWT Configuration**:
   - Set a secret key for JWT signing by adding the following environment variable to the `.env` file:
     ```
     JWT_KEY=<your-secret-key>
     ```
   - Use a strong, random string for the `JWT_KEY` to ensure security.

3. **Example `.env` File**:
   ```env
   DB_USER=workout_user
   DB_NAME=workout_tracker_db
   DB_PASSWORD=securepassword123
   JWT_KEY=your-very-secure-jwt-secret-key
   ```

4. **Install Dependencies**:
   Run the following command to install required Go packages:
   ```bash
   go mod tidy
   ```

5. **Run the API**:
   Start the server using:
   ```bash
   go run main.go
   ```

### Example Usage
#### Create a Workout Plan
```bash
curl -X POST http://localhost:8080/workout-plans/ \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{"ExerciseName": "curlup", "Repititions": 9, "Sets": 3, "Weight": 11}'
```

#### List All Workout Plans
```bash
curl -X GET http://localhost:3000/workouts \
-H "Authorization: Bearer <your-jwt-token>"
```

## Error Handling
The API returns standard HTTP status codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `204 Status No Content`: Successful request with no content
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid JWT
- `500 Internal Server Error`: Server-side error

## Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/Oriseer/workout_tracker/tree/main?tab=MIT-1-ov-file) file for details.

## Contact
For questions or support, reach out to [parianjohnmichael@gmail.com](mailto:parianjohnmichael@gmail.com) or open an issue on GitHub.
