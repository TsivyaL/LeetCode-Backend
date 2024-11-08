# LeetCode-Backend
LeetCode-like System
Overview
This project is a LeetCode-like system designed to allow users to interact with coding challenges. The backend is built with Go (using the Gin framework), and the frontend is developed in Nuxt.js with the Nuxt UI library. Currently, users can create, edit, delete, read, and view questions. The project includes CI/CD pipelines, Docker-based deployment, and utilizes MongoDB as the primary data store.

Features
Question Management: Users can create, view, update, and delete coding questions.
Backend: Built with Go and the Gin framework, providing fast and efficient APIs for managing questions.
Frontend: Developed with Nuxt.js, offering a modern and responsive UI.
Database: MongoDB is used for data storage, accessible via the Go MongoDB driver.
Docker Deployment: Both backend and frontend are containerized using Docker, allowing easy deployment and environment consistency.
CI/CD Integration: Automated testing and deployment pipelines are set up to streamline development.
Project Structure
Backend: Contains the API and logic for managing questions, including routes, controllers, models, and services.
Frontend: Built with Nuxt.js, this directory holds the code for the user interface and interaction with backend APIs.
Docker: Docker configuration for setting up containers for both the backend and frontend.
Endpoints
Questions API
GET /questions - Retrieves all questions.
GET /questions/:id - Retrieves a question by ID.
POST /questions - Creates a new question.
PUT /questions/:id - Updates a question by ID.
DELETE /questions/:id - Deletes a question by ID.
Example Payload for Creating a Question
json
Copy code
{
  "id": "1",
  "title": "What is Go?",
  "body": "Explain what Go is and why it's used.",
  "answer": "Go is a programming language developed by Google.",
  "created": "2024-11-07T15:30:00Z"
}
Setup Instructions
Clone the repository:

bash
Copy code
git clone <repository-url>
Run MongoDB: Ensure MongoDB is running locally on port 27017. Alternatively, configure the MongoDB URI in the backend configuration.

Start the Backend:

Navigate to the Backend directory.
Run the backend server with:
bash
Copy code
go run main.go
Start the Frontend:

Navigate to the Frontend directory.
Run the frontend server with:
bash
Copy code
npm install
npm run dev
Run with Docker:

Build and run the Docker containers for both frontend and backend with:
bash
Copy code
docker-compose up --build
Future Enhancements
Answer Submission and Verification:
Implementing a feature to allow users to submit answers to coding challenges and verify them against the correct solution using Docker to run isolated code containers.
Technologies Used
Backend: Go, Gin framework, MongoDB, Docker
Frontend: Nuxt.js, Nuxt UI library
CI/CD: Docker, GitHub Actions

