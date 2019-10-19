go-skeleton
-----------
This project is used as boilerplate example for the structure of a go microservice. A template to get you started so to speak.

A skeleton for go microservices was the core idea, so that's where the name comes from.


## Motivation
Besides avoiding to copy it from the last microservice I've written, a friend asked how I'd do it.
That's is how I'm doing it right now


## Download

Clone the repository using 

### SSH

```bash
git clone git@github.com:AndersonQ/go-skeleton.git
```

### https

```bash
git clone https://github.com/AndersonQ/go-skeleton.git
```
or using 

```bash
go get -u github.com/AndersonQ/go-skeleton.git
```

and start from there!

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
 