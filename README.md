
### **LeetCode Backend System**  

This project is a backend system designed to work with a LeetCode-like platform. It includes:  
- A backend written in Go (Gin framework) to handle API requests and interact with MongoDB.  
- A MongoDB database for storing questions.  
- A frontend built with Nuxt.js for user interaction.  
- Docker and Kubernetes for development and deployment.  

---

### **System Components**
1. **Backend**:  
   Written in Go using the Gin framework. Handles API requests and database interactions.  

2. **Frontend**:  
   Built with Nuxt.js. Allows users to view questions and submit answers.  

3. **MongoDB**:  
   Stores questions, answers, and related data.  

4. **Docker**:  
   Containerizes the application for simplified development and deployment.  

5. **Kubernetes**:  
   Manages dynamic tasks by running temporary pods to execute user-submitted code.  

---

### **Prerequisites**  
Ensure the following are installed on your local machine:  
1. **Docker**: Install [Docker](https://docs.docker.com/get-docker/).  
2. **Docker Compose**: Install [Docker Compose](https://docs.docker.com/compose/install/).  
3. **Kubernetes**: Set up Kubernetes in Docker Desktop:  
   - Open Docker Desktop.  
   - Go to **Settings > Kubernetes**.  
   - Enable Kubernetes and wait for it to initialize.  

---

### **Getting Started**

#### **1. Clone the repository**  
Clone the project repository to your local machine:  

```bash
git clone https://github.com/TsivyaL/LeetCode-Backend.git
```
#### **2. Set up Kubernetes configuration**  
If you're using Kubernetes for dynamic tasks, follow these steps to configure it:

##### a. Enable Kubernetes in Docker Desktop:
- Open **Docker Desktop**.  
- Go to **Settings > Kubernetes**.  
- Make sure Kubernetes is enabled and running.

##### b. Locate the Kubernetes configuration file:  
The `kubeconfig` file is usually stored under `~/.kube/config`. To verify its location, run:  

```bash
kubectl config view --minify --output "jsonpath={.config}"
```

##### c. Update the `.env` file:  
Add the path to your Kubernetes config file in the `.env` file. Replace `<path_to_kubeconfig>` with the actual file path:

```bash
KUBE_CONFIG_PATH=<path_to_kubeconfig>
```

For example, if using Docker Desktop's Kubernetes, the path might be:  

```bash
KUBE_CONFIG_PATH=/Users/yourusername/.kube/config
```

##### d. Verify Kubernetes setup:  
Ensure Kubernetes is running and pods are configured correctly:

```bash
kubectl get pods
```

---
#### **3. Build and start the application**  
Navigate into the project directory:  

```bash
cd LeetCode-Backend
```

Run the following command to build and start the application using Docker Compose:  

```bash
docker-compose up --build
```



### **System Functionality**
1. The backend manages CRUD operations for questions and handles user submissions.  
2. Kubernetes is used for dynamic task execution by creating temporary pods.  
3. Docker Compose simplifies the local development environment.  

Once all components are up and running, the system will be fully operational and ready to interact with the frontend and backend services.  

---

