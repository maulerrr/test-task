# Project Setup and Running Instructions

## Prerequisites

- Go (version 1.20+)
- PostgreSQL
- Docker (optional, for running PostgreSQL in a container)

## Setting Up the Project

1. **Clone the Repository**

   ```bash
   git clone https://github.com/maulerrr/test-task.git
   cd test-task
   ```

2. **Install Dependencies**

   Make sure you have Go installed and run the following command to download the required Go modules:

   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables**

   Create a `.env` file in the root directory and configure the following environment variables:

   ```env
   PORT=your-app-port
   DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
   JWT_SECRET_KEY=your-jwt-secret-key
   ```

   Replace the placeholder values with your actual configuration.

4. **Set Up the Database**

   If you are not using Docker, make sure PostgreSQL is installed and running on your system.

   If you are using Docker, you can run PostgreSQL with:

   ```bash
   docker run --name your-db-container -e POSTGRES_USER=your-db-username -e POSTGRES_PASSWORD=your-db-password -e POSTGRES_DB=your-db-name -p 5432:5432 -d postgres
   ```

5. **Run the Application**

   You can start the application using:

   ```bash
   go run .\cmd\main.go
   ```

   By default, the application will run on `localhost:8080`.

6. **Accessing Endpoints**

   - **Register User**: `POST /api/v1/auth/register`
   - **Login User**: `POST /api/v1/auth/login`
   - **Refresh Tokens**: `POST /api/v1/auth/refresh-tokens`
   - **Issue Tokens**: `POST /api/v1/auth/issue-tokens/{id}`

   You can use tools like [Postman](https://www.postman.com/) to test these endpoints.

## Testing

To run tests, use:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM for Go
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing

---