# LeetCode Backend System

This project is a backend system designed to work with a LeetCode-like platform. It includes a backend written in Go (Gin framework), a MongoDB database for storing questions, and a frontend built with Nuxt.js. The backend handles the creation, reading, updating, and deletion of questions, while the frontend allows users to interact with the system.

## System Components

- **Backend**: Written in Go using the Gin framework. Handles API requests and interacts with the MongoDB database.
- **Frontend**: Built with Nuxt.js. Displays questions and allows users to submit answers.
- **MongoDB**: Stores questions, answers, and other related data.
- **Docker**: The system uses Docker to containerize the application, simplifying the development and deployment process.

## Prerequisites

Before you begin, ensure that you have the following installed on your local machine:

- **Docker**: [Install Docker](https://www.docker.com/get-started)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

### 1. Clone the repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/TsivyaL/LeetCode-Backend.git

#### 2. Build and start the application
Navigate into the project directory:

cd LeetCode-Backend

docker-compose up --build