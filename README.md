# Microblogging API

This project implements a microblogging API with support for tweets, follows, and timeline management. It features real-time updates through event-driven architecture and performance optimization through caching.

The service is built using Go and separated into distinct layers (handlers, services, repositories). It provides RESTful endpoints for managing tweets and user relationships, with PostgreSQL as the primary data store and Redis for caching frequently accessed data. The development environment is fully containerized using Docker, enabling easy setup and consistent development experience across different environments.

## Repository Structure

```
.
├── cmd/api/                    # Application entry point and server initialization
├── internal/                   # Internal application code
│   ├── config/                # Configuration management and environment settings
│   ├── handler/               # HTTP request handlers and routing logic
│   ├── model/                 # Data models and domain entities
│   ├── repository/            # Data access layer implementations
│   ├── service/               # Business logic and use case implementations
│   └── utils/                 # Shared utilities and helper functions
├── cache/                     # Cache implementation
├── db/                       # Database connection and management
├── docker-compose.yml        # Docker services configuration
├── Dockerfile                # Container image definition
└── testhelpers/             # Testing utilities and helpers
```

## Usage Instructions

### Prerequisites

- Docker and Docker Compose
- Go 1.23 or later
- Redis server (for local development)
- PostgreSQL (for local development)

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd <repository-directory>
```

2. Start the services using Docker Compose:

```bash
docker-compose up -d
```

### Quick Start

1. Configure your environment by copying the example configuration:

```bash
cp internal/config/config_example.yaml internal/config/config_local.yaml
```

2. The API will be available at `http://localhost:5000`

3. Basic usage examples:

Note: user_id is always required in the request url as param and should be an UUUID value.

```bash
# Create a new tweet
curl -X POST http://localhost:5000/api/v1/tweets?user_id=XXXXX \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello, World!", "user_id":XXXX}'

# Follow a user
curl -X POST http://localhost:5000/api/v1/follows?user_id=XXXXX \
  -H "Content-Type: application/json" \
  -d '{"follower_id": "1", "following_id": "2"}'

# Get user timeline
curl http://localhost:5000/api/v1/tweets/timeline?user_id=XXXXX
```

### Troubleshooting

Common Issues:

1. Database Connection Issues

```bash
# Check database container status
docker-compose ps blog-db
# View database logs
docker-compose logs blog-db
```

2. API Service Not Starting

- Verify the database is healthy:

```bash
docker-compose ps
```

- Check API logs:

```bash
docker-compose logs api
```

## Data Flow

The application follows a layered architecture where requests flow through handlers, services, and repositories, with caching at appropriate levels.

Key component interactions:

1. Handlers receive HTTP requests and validate input
2. Services implement business logic and manage transactions
3. Cache layer provides fast access to frequently requested data
4. Repository layer handles data persistence and retrieval
5. Event system enables real-time updates across components
6. Authentication middleware validates requests
7. Error handling occurs at each layer with appropriate responses

## Infrastructure

![Infrastructure diagram](./docs/infra.svg)

### Database

- PostgreSQL instance (`blog-db`)
  - Exposed on port 5432
  - Configured with health checks
  - Persistent volume: `api_postgres_data`

### API Service

- Container: `api`
  - Exposed on port 5000
  - Live reload enabled through Air
  - Depends on database availability

### Testing

- Dedicated test runner service
  - Isolated environment for running tests
  - Mounted source code for real-time updates

## Deployment

### Prerequisites

- Docker and Docker Compose installed
- Available ports: 5000 (API) and 5432 (PostgreSQL)

### Deployment Steps

1. Build and start services:

```bash
docker-compose up --build -d
```

2. Verify deployment:

```bash
docker-compose ps
```

3. Monitor logs:

```bash
docker-compose logs -f
```
