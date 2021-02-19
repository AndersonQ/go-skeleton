go-skeleton
-----------
A skeleton for go microservices

# WIP: opentracing

 1. start _jaeger_ tracer/collector to send the data to with:
````shell
docker run -d \
  --name jaeger \
  --rm -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.21
````
 2. run two instances of this application

```shell
APP_NAME=first_app go run main.go
APP_NAME=second_app PORT=8001 go run main.go
```

 3. run graphql using branch `opentracing`
Run a redis on docker so graphql won't be logging thousands of errors :/ 
```shell
docker run --rm --name graphql-redis -d -p 6379:6379 redis:alpine
npm run start:dev
```

 4. execute the _tracing_ query a few times
```graphql
query {
  tracing
}
```

 5. open [Jaeger UI](http://localhost:16686/search)

 6. have fun :)

## Dependencies

```bash
make deps
```

## Run

### Local
```bash
go run main.go
```

#### Docker
```bash
docker run --rm -e ENV=dev -p8000:8000 go-skeleton
```

## Build

### Local
```bash
make build
```

### Docker

```bash
docker build -t go-skeleton .
```

## Test

```bash
make test
```
 
## Licence
See [LICENCE](LICENSE)
