# Flutter Go CRUD Backend

A high-performance REST API backend built with Go, Fiber, and GORM for managing items with full CRUD operations, pagination, search, and sorting.

## Features

- ✅ **Full CRUD Operations**: Create, Read, Update, Delete items
- 🔍 **Search**: Full-text search across name and description
- 📄 **Pagination**: Efficient pagination with configurable page size
- 🔄 **Sorting**: Sort by name, price, date, or status
- 🗑️ **Soft Delete**: Items can be soft-deleted and restored
- 🔒 **Optimistic Locking**: Version-based concurrency control
- 🚀 **High Performance**: Built with Fiber web framework
- 💾 **Flexible Database**: Supports SQLite and PostgreSQL

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite (default) or PostgreSQL

### Installation

1. **Clone and navigate to backend directory**:
   ```bash
   cd backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Create environment file** (optional):
   ```bash
   cp .env.example .env
   ```

4. **Seed the database** (optional):
   ```bash
   go run cmd/seed/main.go
   ```

5. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080` by default.

## Docker Deployment

### Using Docker

#### Build the Docker Image
```bash
cd backend
docker build -t flutter-go-backend:latest .
```

#### Run with SQLite (Development)
```bash
docker run -d \
  --name flutter-backend \
  -p 8080:8080 \
  -e DB_DRIVER=sqlite \
  -e SQLITE_PATH=/app/data/data.db \
  -v backend-data:/app/data \
  flutter-go-backend:latest
```

#### Run with PostgreSQL
```bash
# First, create a network
docker network create app-network

# Run PostgreSQL
docker run -d \
  --name postgres \
  --network app-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=app \
  -v postgres-data:/var/lib/postgresql/data \
  postgres:15-alpine

# Run the backend
docker run -d \
  --name flutter-backend \
  --network app-network \
  -p 8080:8080 \
  -e DB_DRIVER=postgres \
  -e DB_DSN=postgres://postgres:postgres@postgres:5432/app?sslmode=disable \
  flutter-go-backend:latest
```

### Using Docker Compose (Recommended)

Docker Compose provides a complete development environment with PostgreSQL:

#### Start Services
```bash
cd backend
docker-compose up -d
```

#### View Logs
```bash
docker-compose logs -f backend
```

#### Stop Services
```bash
docker-compose down
```

#### Rebuild After Code Changes
```bash
docker-compose up -d --build
```

#### Clean Up (Remove Volumes)
```bash
docker-compose down -v
```

### Docker Environment Variables

Override any configuration using environment variables:

```bash
docker run -d \
  -p 8080:8080 \
  -e PORT=8080 \
  -e AUTO_MIGRATE=true \
  -e CORS_ORIGINS=http://localhost:*,http://127.0.0.1:* \
  -e DB_DRIVER=sqlite \
  -e SQLITE_PATH=/app/data/data.db \
  -v backend-data:/app/data \
  flutter-go-backend:latest
```

### Seeding Data in Docker

To seed the database in a running container:

```bash
# For SQLite
docker exec -it flutter-backend sh -c "go run cmd/seed/main.go"

# For Docker Compose
docker-compose exec backend sh -c "go run cmd/seed/main.go"
```

Or build a custom image with seed data:

```dockerfile
# Add to Dockerfile
COPY cmd/seed/main.go ./cmd/seed/
RUN CGO_ENABLED=1 go build -o seed ./cmd/seed/main.go
# Then run: docker exec flutter-backend ./seed
```

## Configuration

Configure the application using environment variables (`.env` file or system environment):

### Server Configuration
```env
PORT=8080                    # Server port
AUTO_MIGRATE=true           # Auto-migrate database schema
```

### CORS Configuration
```env
# Comma-separated list of allowed origins (use * for wildcard port)
CORS_ORIGINS=http://localhost:*,http://127.0.0.1:*
```

### Database Configuration

**SQLite (Default)**:
```env
DB_DRIVER=sqlite
SQLITE_PATH=data.db
```

**PostgreSQL**:
```env
DB_DRIVER=postgres
DB_DSN=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## API Endpoints

### Health Check
```
GET /healthz
```

### Items

#### List Items (with pagination, search, sort)
```
GET /api/items?page=1&pageSize=10&search=query&sortBy=name&sortOrder=asc
```

**Query Parameters**:
- `page` (int): Page number (default: 1)
- `pageSize` (int): Items per page (default: 10, max: 100)
- `search` or `q` (string): Search query for name/description
- `sortBy` or `sort` (string): Sort field - `name`, `price`, `createdAt`, `status`
- `sortOrder` or `order` (string): Sort direction - `asc` or `desc`
- `status` (string): Filter by status
- `minPrice` (float): Minimum price filter
- `maxPrice` (float): Maximum price filter
- `includeDeleted` (bool): Include soft-deleted items
- `onlyDeleted` (bool): Show only deleted items

**Response**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Item Name",
      "description": "Item description",
      "price": 99.99,
      "status": "active",
      "version": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 10,
  "total": 100,
  "totalPages": 10
}
```

#### Get Single Item
```
GET /api/items/:id
```

#### Create Item
```
POST /api/items
Content-Type: application/json

{
  "name": "Item Name",
  "description": "Item description",
  "price": 99.99,
  "status": "active"
}
```

#### Update Item
```
PUT /api/items/:id
Content-Type: application/json

{
  "name": "Updated Name",
  "description": "Updated description",
  "price": 149.99,
  "status": "inactive",
  "version": 1
}
```

**Note**: `version` is required for optimistic locking to prevent concurrent update conflicts.

#### Delete Item (Soft Delete)
```
DELETE /api/items/:id
```

#### Bulk Delete Items
```
DELETE /api/items
Content-Type: application/json

{
  "ids": [1, 2, 3],
  "hard": false
}
```

#### Restore Deleted Item
```
POST /api/items/:id/restore
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection & migration
│   ├── dto/
│   │   └── requests.go          # Request DTOs
│   ├── handlers/
│   │   └── item_handler.go     # HTTP handlers
│   ├── models/
│   │   └── item.go              # Data models
│   ├── routes/
│   │   └── routes.go            # Route registration
│   └── utils/
│       └── query.go             # Query parsing utilities
├── .env.example                 # Environment variables example
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o server cmd/server/main.go
./server
```

### Database Migration

The application automatically migrates the database schema on startup when `AUTO_MIGRATE=true`.

To manually migrate:
```go
database.AutoMigrate(db)
```

## Features in Detail

### Soft Delete
Items are soft-deleted by default, meaning they're marked as deleted but not removed from the database. This allows for:
- Data recovery
- Audit trails
- Compliance requirements

Use `hard: true` in bulk delete to permanently remove items.

### Optimistic Locking
The `version` field prevents lost updates in concurrent scenarios:
1. Client reads item with version 1
2. Client modifies and sends update with version 1
3. Server only updates if current version is still 1
4. Version increments to 2 on successful update
5. Concurrent updates with old version fail with 409 Conflict

### Search
Full-text search across `name` and `description` fields using case-insensitive LIKE queries.

### Pagination
Efficient pagination with:
- Configurable page size (max 100)
- Total count
- Total pages calculation
- Offset-based pagination

## Troubleshooting

### Port Already in Use
```bash
# Change port in .env
PORT=3000
```

### Database Connection Issues

**SQLite**:
- Ensure write permissions in the directory
- Check `SQLITE_PATH` is correct

**PostgreSQL**:
- Verify connection string
- Ensure PostgreSQL is running
- Check credentials and database exists

### CORS Issues
Update `CORS_ORIGINS` in `.env` to include your frontend URL:
```env
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
```

## License

MIT
