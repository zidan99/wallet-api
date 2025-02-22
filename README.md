## Tech Stacks

- Go v1.20.5
- PostgreSQL v14.\*

## Installation & Setup

### 1. Clone project

```bash
# Using HTTPS
git clone https://github.com/zidan99/wallet-api.git
```

```bash
# Using SSH
git clone git@github.com:zidan99/wallet-api.git
```

### 2. Copy and paste `.env.example` and rename it to `.env`

```bash
cp .env.example .env
```

After that, setup the required field like: app config, database connection, logging config

```text
APP_NAME="APP_NAME"
APP_PORT="3000"
APP_ENV="development"

# PostgreSQL Connection
DB_HOST="localhost"
DB_PORT="5432"
DB_USERNAME="username"
DB_PASSWORD="password"
DB_DATABASE="database"

# Log Environment
LOG_PATH=
LOG_MODE=
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the application

```bash
go run main.go
```

Done. You're ready to rock!

## Database Migrations & Seeders

### 1. Install golang migrate for migration and seeders

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. Create new migration or seeder

```bash
migrate create -ext sql -dir database/postgresql/migrations your_migration_name
```

```bash
migrate create -ext sql -dir database/postgresql/migrations your_seeder_name
```

### 3. Run database migrations

```bash
# Run migrations
migrate -database ${DB_MIGRATION_CONNECTION} -path database/postgresql/migrations up
```

```bash
# Rollback migrations
migrate -database ${DB_MIGRATION_CONNECTION} -path database/postgresql/migrations down
```