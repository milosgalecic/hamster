# Hamster - Trucking and Logistics Dashboard

Hamster is a Go-based web application designed for trucking and logistics companies to track jobs, drivers, and trucks efficiently. This dashboard provides a central interface for managing job assignments, monitoring driver statuses, and managing truck resources.

## Running the Application

To run the Hamster web application locally:

1. Ensure you have Go installed on your system.
2. Install MySql and follow the instructions in db_backup.sql to recreate the DB.
3. Clone the repository and run the app:
   ```bash
   git clone https://github.com/milosgalecic/hamster.git
   cd hamster
   go mod tidy
   go run ./cmd/web

## Dependencies

This project uses the following Go libraries:

- **github.com/go-sql-driver/mysql**: MySQL driver for Go.
  
- **github.com/julienschmidt/httprouter**: A fast and efficient HTTP router for Go.
  
- **github.com/justinas/alice**: A lightweight middleware stack for Go that allows easy composition of multiple middleware functions.
  
- **filippo.io/edwards25519**: An implementation of the Ed25519 signature scheme, which is used for secure authentication or cryptographic signing in the application.

- **github.com/alexedwards/scs/mysqlstore**: A MySQL-backed session store for the `scs` session management library. This allows session data (like logged-in users) to be stored persistently in the MySQL database.

- **github.com/alexedwards/scs/v2**: A session management package for Go. It provides a simple interface for managing user sessions and supports features like session expiration, secure cookies, and more.
