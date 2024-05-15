# 🚀 GoLang API with PostgreSQL

This project is aimed at creating a simple yet robust API using GoLang, which supports common `Create`, `Read`, `Update` and `Delete (CRUD)` operations on a PostgreSQL database. It provides a solid foundation for improving and understanding GoLang and its efficacy in creating a server side application.

## 📁 Project Structure

Here's a breakdown of the main components of the project:

- `Docker Compose:` File for setting up the PostgreSQL database.
- `db.go:` This file contains necessary setup and functions for initial connection with PostgreSQL database.
- `handle.go:` This is where all handler functions live. These are the functions that execute instructions as per the API requests.
- `model.go:` This file describes struct types and relevant functions.
- `sql.go:` This file contains functions for generating CRUD SQL scripts.
- `main.go:` The controlling file of the application. It is where the router and related handlers are defined.
- `DB_DDL.sql:` File for Data Definition Language (DDL) script and trigger function for automatic updates of 'updated_at' timestamps.
- `SAMPLE_DATA.sql:` Contains a set of sample data for testing.

## 🚀 Getting Started

1. Clone the repository to your local machine

```bash
git clone https://github.com/billybillysss/golang-api.git
```

2. Navigate to the cloned directory

3. If you have Docker installed, you can build and run the PostgreSQL database container with:

```bash
docker-compose up -d
```

4. Once the Docker container is up and running, execute `main.go`:

```bash
go run .
```

Your application should now be running and ready to accept requests!

## 🧪 Interacting with the API

Once your application is running, you can make CRUD operations via HTTP requests to `localhost: portNumber/path`

The path and its function are as follows:

- POST `/members`: Creates a new record
- GET `/members`/`/members/{id}`: Fetches records
- PUT `/members`/`/members/{id}`: Updates an existing record
- DELETE `/members/{id}`: Deletes a record


