# URLShortener Backend Implementation

## Description
This is the backend GO implementation for the URLShortener API. There is already a FE implementation, you can read more about it [here](../frontend/README.md).

The goal of the app as you have probably guessed is to shorten given urls and redirect user back to the original url when later called.

## Architecture

The project follows [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/The-Clean-Architecture.html) to get:

* **easy unit-testing & mocking**  
* **clear separation of concerns**  
* **plug-and-play transports or providers** 

Key libraries & tools:

| Concern | Library |
|---------|---------|
| HTTP server | [Echo](https://echo.labstack.com/) |
| DB migrations | [golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| Testing | [Ginkgo](https://onsi.github.io/ginkgo/) + [Gomega](https://onsi.github.io/gomega/) |
| Fakes | [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) |
| DB | [PostgreSQL] (https://www.postgresql.org/) |
| Caching | [Redis] (https://github.com/redis/go-redis) |

## How to run the server

To run the server using docker:
```bash
make run-docker
```

## How to run the tests

To run the unit tests:

```bash
make test
```

## Example requests

### To shorten a URL

```bash
curl -X POST http://127.0.0.1:5000/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
        "url": "https://www.example.com"
      }'
```

### To test a shortened url

```bash
curl -X GET http://127.0.0.1:5000/abcd1234
```
