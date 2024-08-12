FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/local/src/users-service

COPY . .

RUN apk add --no-cache bash

RUN go mod download


RUN go build -o users-service ./cmd/users-service/main.go


FROM alpine AS runner

WORKDIR /users-service

RUN apk add --no-cache bash

COPY --from=builder /usr/local/src/users-service/users-service .

CMD ["./users-service"]




