# Server backend for Cup Predictor written in Go

# Requirements

-   Go or Docker
-   Postgres server with user `postgres/postgrespwd` and empty database `db`

# Running

```
go get .
go run .
```

# Docker

```
docker build  . -t <image name>
docker run -p 7000:7000 -e POSTGRES_URL=<postgres url> -e JWT_KEY=<secret key> <image name>
```
