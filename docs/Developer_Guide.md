<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer Guide for Ubuntu Autoinstall Webhook](#developer-guide-for-ubuntu-autoinstall-webhook)
  - [**1. Introduction**](#1-introduction)
  - [**2. Setting Up the Development Environment**](#2-setting-up-the-development-environment)
    - [**2.1. Prerequisites**](#21-prerequisites)
    - [**2.2. Cloning the Repository**](#22-cloning-the-repository)
    - [**2.3. Installing Dependencies**](#23-installing-dependencies)
    - [**2.4. Running the Webhook Locally**](#24-running-the-webhook-locally)
  - [**3. Webhook API Documentation**](#3-webhook-api-documentation)
    - [**3.1. Endpoints**](#31-endpoints)
      - [`POST /hardware-info`](#post-hardware-info)
      - [`POST /logs`](#post-logs)
      - [`POST /status`](#post-status)
  - [**4. Best Practices**](#4-best-practices)
  - [**5. Deployment Instructions**](#5-deployment-instructions)
    - [**5.1. Docker Deployment**](#51-docker-deployment)
    - [**5.2. Kubernetes Deployment**](#52-kubernetes-deployment)
  - [**6. Conclusion**](#6-conclusion)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide for Ubuntu Autoinstall Webhook

## **1. Introduction**

This guide provides developers with information on setting up, modifying, and extending the Ubuntu Autoinstall Webhook. It includes setup instructions, API documentation, and best practices for development and deployment.

## **2. Setting Up the Development Environment**

### **2.1. Prerequisites**

- Docker & Docker Compose
- Python 3.8+
- PostgreSQL or SQLite (for local development)
- `pipenv` or `virtualenv` for dependency management
- Git

### **2.2. Cloning the Repository**

```bash
$ git clone https://github.com/example/autoinstall-webhook.git
$ cd autoinstall-webhook
```

### **2.3. Installing Dependencies**

```bash
$ pip install -r requirements.txt
```

### **2.4. Running the Webhook Locally**

```bash
$ python app.py
```

## **3. Webhook API Documentation**

### **3.1. Endpoints**

#### `POST /hardware-info`

- **Description**: Receives and stores client hardware details.
- **Request Body (JSON):**
```json
{
  "mac_address": "AA:BB:CC:DD:EE:FF",
  "cpu": "Intel i7",
  "ram": "16GB",
  "disk": "500GB NVMe"
}
```
- **Response:** HTTP 200 on success.

#### `POST /logs`

- **Description**: Receives installation logs from clients.
- **Request Body (Text/Log Format)**
- **Response:** HTTP 200 on success.

#### `POST /status`

- **Description**: Updates installation progress.
- **Request Body (JSON):**
```json
{
  "mac_address": "AA:BB:CC:DD:EE:FF",
  "status": "OS Installed"
}
```
- **Response:** HTTP 200 on success.

## **4. Best Practices**

- Use structured logging to ensure troubleshooting is easier.
- Follow RESTful API principles for any modifications.
- Write unit tests for new features before committing changes.
- Use environment variables for configuration to support different deployment environments.

## **5. Deployment Instructions**

### **5.1. Docker Deployment**

```bash
$ docker build -t autoinstall-webhook .
$ docker run -d -p 5000:5000 autoinstall-webhook
```

### **5.2. Kubernetes Deployment**

- Define a Kubernetes deployment file.
- Ensure persistent storage for logs and database.

## **6. Conclusion**

This guide provides a complete overview of setting up and managing the webhook service for Ubuntu Autoinstall. Developers should follow best practices and API guidelines to maintain system reliability and scalability.
