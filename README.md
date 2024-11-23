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
- **Kubernetes**: Ensure that Kubernetes is set up in your Docker environment. If you haven't installed it yet, you can enable Kubernetes in Docker Desktop as follows:
  1. Open **Docker Desktop**.
  2. Go to **Settings** > **Kubernetes**.
  3. Enable Kubernetes and wait for it to initialize.
  
## Getting Started

### 1. Clone the repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/TsivyaL/LeetCode-Backend.git
2. Build and start the application
Navigate into the project directory:

bash
Copy code
cd LeetCode-Backend
Run the following command to build and start the application using Docker Compose:

bash
Copy code
docker-compose up --build
3. Set up Kubernetes configuration
If you are using Kubernetes to run your dynamic tasks, follow these steps to configure the Kubernetes connection:

Enable Kubernetes in Docker Desktop:

Go to Docker Desktop > Settings > Kubernetes.
Make sure Kubernetes is enabled and running.
Find the Kubernetes config file: The Kubernetes configuration file (kubeconfig) is usually stored in your user home directory under ~/.kube/config. To verify its location within Docker, you can run the following command inside the Docker Desktop terminal:

bash
Copy code
kubectl config view --minify --output "jsonpath={.config}"
This will show the path to your kubeconfig file.

Update your .env file: Once you've located the kubeconfig file, update your .env file with the correct path to the configuration file.

Open the .env file in your project directory and add the following line (replace <path_to_kubeconfig> with the actual path):

bash
KUBE_CONFIG_PATH=<path_to_kubeconfig>
For example, if you're running Kubernetes locally through Docker Desktop, the path might look like this:

bash
KUBE_CONFIG_PATH=/Users/yourusername/.kube/config
Verify Kubernetes setup: Ensure that Kubernetes is running correctly by checking the pods:

bash
kubectl get pods
If everything is set up correctly, the system will be ready to interact with Kubernetes for executing dynamic tasks.
