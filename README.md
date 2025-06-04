# SequenceSender

A Go application for creating and managing email sequence campaigns. This application uses handler-service 
design pattern that provides good separation of concern without introducing too much boilerplate code and compromising
extensibility.

## Prerequisites

- Go 1.21 or higher
- Docker/OrbStack and Docker Compose
- PostgreSQL

## Setup

### 1. Environment Configuration

Copy the `placeholder.env` file to `.env` and update the values:

```bash
cp placeholder.env .env
```

Update the `.env` file with your environment settings

### 2. Setup local database

Use composeer file to setup the PostgreSQL database. 

```bash
docker-compose up -d
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Build the Application

Building the API. Use development flag only in development.

```bash
go build -o bin/api/main cmd/api/main.go
bin/api/main --development
```

## Database Entities

- **mailboxes**: Available email mailboxes
- **sequences**: Email sequences 
- **sequence_steps**: steps under a sequence
- **contacts**: Contact (targets) information
- **campaigns**: Campaign instances representing runs of a sequence

## Development

### Running PostgreSQL

**Custom Function Option:**
```bash
# Start PostgreSQL with custom UUID v6 function
docker-compose up -d

# Stop PostgreSQL
docker-compose down
```

### Migrations

Migration files can be found in the `migrations/versions` directory. See [golang-migrate](https://github.com/golang-migrate/migrate) for more details. The migration CLI utility can be accessed using 

```sh
go build -o bin/migrator/main cmd/migrator/main.go
sudo chmod +x bin/migrator/migrate 
bin/migrator/migrate --help
```

### Unit Tests 
This project uses [vektra/mockery](https://github.com/vektra/mockery) for generating mocks and 
[stretchr/testify](https://github.com/stretchr/testify) for assert helpers. 

See .mockery.yaml for mockery configurations.

Run 
```sh
mockery
```

or 
```sh
make mocks
```

to generate mocks based on .mockery.yaml.

## TODO 
- ~~Makefile aliases~~
- ~~Unit tests - test coverage Makefile alias~~
- [Bonus] 
  - Add utility to scaffold handler/service/storage methods 
  - Add di management library like uber/wire.
  - Configure linters
  - Error wrapping, propogation and central logging
  - Validator library

## Order of imports
- standard library 
- internal package dependencies 
- external libraries