FROM golang:1.17-alpine AS base
WORKDIR /app

RUN apk update && apk add --no-cache gcc
RUN apk add libc-dev

RUN go get github.com/rubenv/sql-migrate/...

COPY go.mod ./

# System dependencies
RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

# Hot reloading mod
RUN go get -u github.com/cosmtrek/air
EXPOSE 8080

ENTRYPOINT ["air"]

RUN go mod download
RUN apk add build-base

COPY . .

#### Executable builder
#FROM base AS builder
#WORKDIR /usr/src/app
#
## Application dependencies
#COPY ./ /usr/src/app
#RUN go mod download \
#    && go mod verify
#
#RUN apk add build-base
#
#COPY . .
#
#RUN go build -o /usr/src/app -a .

#### Production
#FROM alpine:latest
#
#RUN apk update \
#    && apk add --no-cache \
#    ca-certificates \
#    curl \
#    tzdata \
#    && update-ca-certificates
#
## Copy executable
#COPY --from=builder /app/social-network /usr/local/bin/social-network
#EXPOSE 8080
#
#ENTRYPOINT ["/usr/local/bin/social-network"]
