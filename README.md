# TestApp

TestApp is a small service for reading GitHub repository information and managing repository subscriptions. The public interface is HTTP; internal services communicate through gRPC and Kafka, and use PostgreSQL for persistence.

## Architecture

```mermaid
flowchart LR
    Client[HTTP client] --> Gateway[Gateway]

    Gateway -->|gRPC| Processor[Processor]
    Gateway -->|gRPC| Subscriber[Subscriber]

    Processor -->|SQL| ProcessorDB[(PostgreSQL)]
    Processor -->|task request| Kafka[(Kafka)]
    Kafka -->|task response| Processor

    Kafka -->|task request| Collector[Collector]
    Collector -->|task response| Kafka

    Processor -->|gRPC| Subscriber
    Collector -->|gRPC| Subscriber

    Collector -->|fetch repositories| GitHub[GitHub API]
    Subscriber -->|validate repository| GitHub
    Subscriber -->|SQL| SubscriberDB[(PostgreSQL)]

    MigrateSubscriber[subscriber migrations] --> SubscriberDB
    MigrateProcessor[processor migrations] --> ProcessorDB
```

- **Gateway** exposes the HTTP API and Swagger UI.
- **Processor** caches repository data. If the data is missing, it requests it through Kafka.
- **Collector** reads repository data from the GitHub API, publishes task responses to Kafka, and periodically refreshes subscribed repositories.
- **Subscriber** creates, deletes, and lists repository subscriptions.
- **PostgreSQL** stores subscriptions and cached repository data.
- **Kafka** decouples repository fetch requests from HTTP/gRPC request handling.

## Requirements

- Docker
- Docker Compose
- GitHub token for authenticated GitHub API requests

Create a local `.env` file:

```bash
cp .env.example .env
```

Then set:

```env
GITHUB_TOKEN=github_pat_...
```

The `.env` file is used by Docker Compose for variable substitution.

## Running

Start the whole stack:

```bash
docker compose up --build
```

After startup:

- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/docs/swagger/index.html`
- Kafka from host: `localhost:9094`
- Subscriber PostgreSQL from host: `localhost:5432`
- Processor PostgreSQL from host: `localhost:5433`

Stop services:

```bash
docker compose down
```

Stop services and remove database volumes:

```bash
docker compose down -v
```

## HTTP API

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/repositories/info?url=https://github.com/{owner}/{repo}` | Get repository information |
| `POST` | `/api/subscriptions` | Create a subscription |
| `DELETE` | `/api/subscriptions/{owner}/{repo}` | Delete a subscription |
| `GET` | `/api/subscriptions` | List subscriptions |
| `GET` | `/api/subscriptions/info` | Get information about subscribed repositories |
| `GET` | `/api/ping` | Check service health |

Repository information can return:

- `200 OK` when repository data is already available.
- `202 Accepted` when the repository fetch task has been queued and data is still being prepared.
- `400`, `404`, `502`, or `500` for invalid input, missing repositories, GitHub unavailability, or internal errors.

Create a subscription:

```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H 'Content-Type: application/json' \
  -d '{"owner":"octocat","repo":"Hello-World"}'
```

Get repository information:

```bash
curl 'http://localhost:8080/api/repositories/info?url=https://github.com/octocat/Hello-World'
```

Check service health:

```bash
curl http://localhost:8080/api/ping
```

## Local Service Ports

| Service | Container port | Host port |
|---|---:|---:|
| Gateway HTTP | `8080` | `8080` |
| Processor gRPC | `50051` | `50051` |
| Collector gRPC | `50052` | `50052` |
| Subscriber gRPC | `50053` | `50053` |
| Kafka external listener | `9094` | `9094` |
| Subscriber PostgreSQL | `5432` | `5432` |
| Processor PostgreSQL | `5432` | `5433` |

## Development

Run Go tests:

```bash
go test ./...
```

Regenerate subscriber sqlc code:

```bash
cd subscriber/internal/adapter/postgres
sqlc generate
```

Regenerate processor sqlc code:

```bash
cd processor/internal/adapter/postgres
sqlc generate
```

Regenerate Swagger documentation:

```bash
swag init \
  -g main.go \
  -d gateway/cmd,gateway/internal/controller/http,gateway/internal/domain \
  -o docs
```
