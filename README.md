
# GO Skeleton

> A skeleton for go microservices

## Motivation
> Besides avoiding to copy it from the last microservice I've written, a friend asked how I'd do it.
> That's is how I'm doing it right now


## Dependencies

```bash
make deps
```

## Run

-- Local
```bash
go run main.go
```

-- Docker
```bash
docker run --rm -e ENV=dev -p8000:8000 go-skeleton
```

## Build

-- Local
```bash
make build
```

-- Docker

```bash
docker build -t go-skeleton .
```

## Test

```bash
make test
```
 
## Licence
> See [LICENCE](LICENSE)
