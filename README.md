# book-shop

An example API for a book shop.

## Installation

- [Golang 1.16.5](https://golang.org/doc/install)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [golang-migrate](https://github.com/golang-migrate/migrate) - `go get -tags 'postgres' -u https://github.com/golang-migrate/migrate`

#### Setting Up Your Development Environment

Create docker container for dependencies:

```
docker-compose up -d
```

Create `book-shop` database:

```
make db_create
```

Migrate necessary tables using `golang-migrate`:

```
make db_migrate
```

## Testing

`Coming Soon`

## API Documentation

`Coming Soon`
