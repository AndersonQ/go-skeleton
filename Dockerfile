# Prepare
FROM golang:1.13-alpine as baseimg

RUN apk --no-cache upgrade && apk --no-cache add git make

# First only download the dependencies, so thid layer can be cahced before we copy the code
COPY ./go.mod ./go.sum ./Makefile /app/
WORKDIR /app/
RUN make deps

# Build
FROM baseimg as builder

COPY . ./
RUN make build

# Run
FROM alpine

COPY --from=builder /app/go-skeleton /opt/
WORKDIR /opt/
ARG ENV
EXPOSE 8000
CMD ["./go-skeleton"]