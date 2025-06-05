
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
  **Request Body**: JSON with workout plan details (e.g., name, exercises, duration).  
  **Response**: JSON with created workout plan details or error message.

- **PUT /workout-plans/**  
  Update an existing workout plan.  
  **Requires Authentication**: Yes  
  **Request Body**: JSON with updated workout plan details.  
  **Response**: JSON with updated workout plan details or error message.

- **DELETE /workout-plans/**  
  Delete a workout plan.  
  **Requires Authentication**: Yes  
  **Request Body**: JSON with workout plan ID.  
  **Response**: Success message or error message.

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
- Node.js (or your preferred runtime environment)
- A database (e.g., MongoDB, PostgreSQL)
- JWT library for your programming language
- HTTP client (e.g., Postman, curl) for testing

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/workout-tracker-api.git
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Set up environment variables (e.g., `.env` file):
   ```bash
   DATABASE_URL=your-database-url
   JWT_SECRET=your-jwt-secret
   PORT=3000
   ```
4. Start the server:
   ```bash
   npm start
   ```

### Example Usage
#### Create a Workout Plan
```bash
curl -X POST http://localhost:3000/workout-plans/ \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{"name": "Morning Routine", "exercises": ["Push-ups", "Squats"], "duration": 30}'
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
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact
For questions or support, reach out to [your-email@example.com](mailto:your-email@example.com) or open an issue on GitHub.

