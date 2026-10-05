

FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=1 GOOS=linux go build -o bandplan ./cmd/server 

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bandplan ./bandplan
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/sql ./sql

RUN addgroup -S -g 10001 bandplan \
    && adduser -S -D -H -u 10001 -G bandplan bandplan \
    && mkdir -p /app/static/uploads/song-images \
    && chown -R bandplan:bandplan /app/static/uploads/song-images

USER 10001:10001

EXPOSE 8080

CMD ["./bandplan"]
