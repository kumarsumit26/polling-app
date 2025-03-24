# Go API Project

This project is a simple API built with Go that interacts with a PostgreSQL database to fetch records. 

## Project Structure

```
go-api-project
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── api
│   │   ├── handler.go   # API request handlers
│   │   └── router.go    # HTTP routes setup
│   ├── db
│   │   ├── postgres.go   # PostgreSQL database connection management
│   │   └── migrations
│   │       └── 0001_initial.sql # Database schema initialization
│   └── models
│       └── record.go    # Data model for records
├── go.mod                # Module definition and dependencies
├── go.sum                # Dependency checksums
└── README.md             # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-api-project
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Set up PostgreSQL:**
   - Ensure you have PostgreSQL installed and running.
   - Create a database for the project.

4. **Run migrations:**
   - Execute the SQL commands in `internal/db/migrations/0001_initial.sql` to set up the initial schema.

5. **Run the application:**
   ```
   go run cmd/main.go
   ```

## API Usage

- **GET /records**
  - Fetches all records from the database.
  - Response: JSON array of records.

## License

This project is licensed under the MIT License.