# E-Commerce and Stock Management System Documentation

This project outlines developing an e-commerce system comprising various services for user authentication, product management, order processing, and warehouse stock management. The system is designed to be scalable, robust, and maintainable, leveraging best practices in software development. This is the solution for implementing the microservices in one repository or we can call it _monorepo_.

## Prerequisite
There are some library need to be installed to develop this project:
1. [Go](https://golang.org/doc/install)
2. [mockgen](https://github.com/uber-go/mock)
3. [GNU Make](https://www.gnu.org/software/make/)
4. [Docker Compose](https://docs.docker.com/compose/install/)
5. [SQLite](https://www.sqlite.org/)

## Table of contents
1. [Service Overview](#services-overview)
2. [Reproduce The Project](#reproduce-the-project)
3. [Postman Collection](#postman-collection)
4. [Project Structure](#project-structure)
5. [Notes](#notes)

## Services Overview
### 1. User Service
- **Authentication:** Implements simple authentication for users to log in using either phone or email.

### 2. Product Service
- **List Products:** Provides an API to retrieve a list of products along with their stock availability from the database.

### 3. Order Service
- **Checkout and Stock Deduction:** Processes customer orders by reserving (locking) stock for ordered products. Ensures stock availability before confirming an order to prevent overselling.
- **Release Stock:** Releases reserved stock if payment is not completed within a specified time frame (e.g., N minutes) using background jobs or timers.

### 4. Shop Service
- **Warehouse Management:** Tracks the association of one or more warehouses with a shop.

### 5. Warehouse Service
- **Stock Management:** Handles inventory levels and updates.
- **Transfer Products:** Allows product stock transfer between warehouses. Updates stock levels accordingly.
- **Active/Inactive Warehouses:** Maintains the status of each warehouse. Excludes stock from inactive warehouses from the available stock pool. Provides mechanisms to activate or deactivate warehouses.

## Reproduce The Project
Clone the project
```
git clone https://github.com/DangerZombie/edot-ecommerce.git
```

Ensure module
```
go mod tidy
```

## Postman Collection
Use the Postman Collection for e2e testing. If you need the Postman Collection, please contact me. 😄

## Project Structure
```
project/ 
├── data/ 
│   ├── ecommerece.db/ 
│   └── ...
├── micro-services/ 
│   ├── order/ 
│   │   ├── config/
│   │   ├── cron/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── migrations/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── main.go 
│   ├── product/ 
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── migrations/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── main.go 
│   ├── shop/ 
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── migrations/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── main.go 
│   ├── user/ 
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── migrations/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── main.go 
│   ├── warehouse/ 
│   │   ├── config/
│   │   ├── cron/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── migrations/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── main.go 
│   └── ...
├── go.mod 
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── ...
```

## Notes
- This project is developed using Go version 1.22.0
- This project structure represents the microservices approach with simplification using _monorepo_
- Those services running on different port from 7001 - 7005