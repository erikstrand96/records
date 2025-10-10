FROM golang:1.25.2-alpine AS builder
WORKDIR /usr/src/app


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /usr/local/bin/records .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1000 appuser

WORKDIR /home/appuser

COPY --from=builder /usr/local/bin/records .

USER appuser

EXPOSE 7002

CMD ["./records"]