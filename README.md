# EzyMark

EzyMark is an e-commerce system that provides APIs for managing products, orders, and users. This project is built with Go and uses the Fiber framework. EzyMark is designed to be highly concurrent and scalable, leveraging modern technologies to ensure high performance and reliability.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technical Stack](#technical-stack)
- [System Architecture](#system-architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [License](#license)

## Introduction

EzyMark is a powerful and flexible e-commerce system that allows users to efficiently manage products, orders, and users. This project is designed to be easily extendable and maintainable, with a focus on high concurrency and scalability.

## Features

Updating...

## Technical Stack

- Backend building blocks

  - [Go](https://github.com/golang/go): The main programming language used for building the application, known for its performance and concurrency support.
  - [Fiber](https://github.com/gofiber/fiber): An Express-inspired web framework for Go, optimized for performance.
  - [Asynq](https://github.com/hibiken/asynq): A Go library for processing asynchronous tasks.
  - [Zerolog](https://github.com/rs/zerolog): A fast and lightweight logging library for Go.
  - [Sonic](https://github.com/bytedance/sonic): A fast JSON library for Go.
  - [JWT](https://github.com/golang-jwt/jwt): A library for working with JSON Web Tokens.
  - [Viper](https://github.com/spf13/viper): A complete configuration solution for Go applications.
  - [Goose](https://github.com/pressly/goose): A database migration tool for managing schema changes.
  - [Sqlc](https://github.com/sqlc-dev/sqlc): A tool for generating type-safe Go code from SQL queries.

- Infrastructure

  - **PostgreSQL**: A robust, open-source relational database system used for storing persistent data.
  - **Redis**: An in-memory data store used as a database, cache, and message broker to improve performance.
  - **Docker and docker-compose**: Tools for containerizing applications and managing multi-container deployments.
  - **MinIO**: A high-performance, S3-compatible object storage system for storing user-uploaded files.

## System Architecture

Updating...

## Installation

### Prerequisites

Updating...

### Installation Steps

Updating...

## Usage

### API Documentation

API documentation can be accessed at `/swagger/*`.

## Project Structure

```plaintext
ezy-mark/
├── cmd/
│   └── server/
│       └── main.go
├── docs/
├── environment/
├── internal/
│   ├── cache/
│   ├── controllers/
│   ├── consts/
│   ├── db/
│   │   ├── migrations/
│   │   └── query/
│   ├── global/
│   ├── helpers/
│   ├── initialize/
│   ├── mapper/
│   ├── models/
│   ├── middlewares/
│   ├── pkg/
│   ├── response/
│   ├── routes/
│   ├── services/
│   │   └── impl/
│   └── worker/
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

