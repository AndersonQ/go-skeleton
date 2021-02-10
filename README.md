go-skeleton
-----------
A skeleton for go microservices



# WIP: opentracing

 - start _jaeger_ tracer/collector to send the data to with:
````shell
docker run -d --name jaeger \
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
 - run this application

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
