# carddeck
Go REST API for playing cards

## Prerequisites
- [Go](https://go.dev/dl/) (1.22.X). You may use [gvm](https://github.com/moovweb/gvm) to manage multiple go version.
- [Docker](https://docs.docker.com/engine/install/), for spinning up dependencies
- make (optional, most unix already have this)
- [Postgres](https://www.postgresql.org/download/) (alternatively if you want to install the dependencies in local)

## How to Run (Unix)

### Build the application
Executing following command to build the application

```
go build . -o carddeck

(alternatively)
make build
```

### Copy .env.sample to .env file
Copy `.env.sample` file to `.env` file, change the value if needed.
```
cp .env.sample .env
```

### Spin up dependencies
Run docker compose to spin up dependencies
```
docker compose up -d
```

### Run the application
Execute following command to run the application. The application will be available in `localhost:8080` by default.
```
.\carddeck

(alternatively#1) 
make run

(alternatively#2, if you don't want to build the app) 
go run . server
```

### Shutting down dependencies
```
docker compose down
```
