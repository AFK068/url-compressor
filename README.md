# Simple URL Compressor

![Coverage](https://img.shields.io/badge/coverage-34%25-brightgreen.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

This is a simple backend project for URL shortening. It provides an easy-to-use API that allows you to save URLs, generate shortened links, and retrieve the original URLs using the shortened ones.

## How to Run

Ensure you have [Docker](https://docs.docker.com/install/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
### Running the Application

You can run the application using one of the following storage options:

#### 1. In-Memory Storage
This option uses in-memory storage. Run the following command:
```bash
STORAGE_TYPE="inmemory" docker-compose up -d
```

#### 2. PostgreSQL Storage
For use PostgreSQL:
```bash
POSTGRES_PASSWORD=<your_password> STORAGE_TYPE="postgres" docker-compose up -d
```
### URL Endpoints

- `GET /url?short-link=` - Retrieves the original URL associated with the provided short link.

- `POST /url` - Creates a new short link. The request body should include:
    ```json
    {
        "url": "http://example.com"
    }
    ```
