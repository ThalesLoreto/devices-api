# Device API

A RESTful API for managing device resources, built with Go following SOLID principles and software engineering best practices.

## Features

- **CRUD Operations**: Create, read, update, and delete device resources
- **Filtering**: Fetch devices by brand or state
- **Business Rules**: Enforces domain validations (e.g., in-use devices cannot be deleted)
- **Database Persistence**: PostgreSQL database with automatic migrations
- **Containerization**: Docker support for easy deployment
- **Comprehensive Testing**: Unit tests with good coverage
- **API Documentation**: Well-documented REST endpoints

## Architecture

The application follows a layered architecture implementing SOLID principles:

### Layers

1. **Presentation Layer** (`internal/handler/`): HTTP handlers for REST API endpoints
2. **Service Layer** (`internal/service/`): Business logic and domain validations
3. **Data Access Layer** (`internal/repository/`): Database operations and data persistence
4. **Domain Layer** (`internal/models/`): Core business entities and domain logic

### Key Components

- **Models**: Device entity with business rules and validations
- **Repository**: Interface-based data access with PostgreSQL implementation
- **Service**: Business logic orchestration and validation
- **Handler**: HTTP request/response handling and routing
- **Config**: Environment-based configuration management
- **Database**: Connection management and migrations

## Device Domain

### Properties
- **ID**: Unique identifier (UUID)
- **Name**: Device name
- **Brand**: Device brand/manufacturer
- **State**: Device state (available, in-use, inactive)
- **Creation Time**: Timestamp when device was created

### Business Rules
- Creation time cannot be updated
- Name and brand cannot be updated if device is in use
- In-use devices cannot be deleted
- All fields except creation time are required for creation

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd device-api
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **View logs**
   ```bash
   docker-compose logs -f api
   ```

5. **Stop services**
   ```bash
   docker-compose down
   ```

The API will be available at `http://localhost:8080`

## Testing

### Run Unit Tests
```bash
go test ./test/... -v
```

### Run Tests with Coverage
```bash
go test ./test/... -v -cover
```

## Configuration

The application uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `deviceapi` | Database name |
| `DB_SSLMODE` | `disable` | Database SSL mode |

## Database Schema

```sql
CREATE TABLE devices (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    state VARCHAR(50) NOT NULL CHECK (state IN ('available', 'in-use', 'inactive')),
    creation_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_devices_brand ON devices(brand);
CREATE INDEX idx_devices_state ON devices(state);
CREATE INDEX idx_devices_creation_time ON devices(creation_time);
```

## Future Improvements

### TODO
