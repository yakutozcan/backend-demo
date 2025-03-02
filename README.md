# backend-demo

### Description

### Overview

In this project, you are tasked with developing a system that automatically sends 2
messages retrieved from the database, which have not yet been sent, every 2 minutes, as
illustrated in the CURL request provided below.

API View Url: https://webhook.site/#!/view/f0c1ed93-8a60-41d8-a661-7b9ce22774d7/

### Prerequisites

- Go 1.23.3
- Docker and Docker Compose
- PostgreSQL
- golang-migrate

### Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Container**: Docker
- **Config**: Viper
- **Logger**: Zap

### Project Structure

```lua
├── app
│ └── api/healthcheck
│── └── health.go
├── config
│ └── config.yaml -- App Configuration File
├── pkg
│ └── config
│── └── config.go
│ └── log
│── └── log.go
├── Dockerfile
├── docker-compose.yml
├── main.go
```

### Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/yakutozcan/backend-demo.git
    ```

2. Start the required services using Docker Compose:

    ```bash
    docker-compose up -d
    ```

3. Run the migration:
   ```bash
      migrate -database "postgres://postgres:password@localhost:5432/backenddemo?sslmode=disable" -path db/migrations up    ```
   ``` 
   ```bash
      migrate -database "postgres://postgres:password@localhost:5432/backenddemo?sslmode=disable" -path db/migrations drop 
   ```

4. Run the application:

    ```bash
    go run .
    ```