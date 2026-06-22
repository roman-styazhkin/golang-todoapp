FROM golang:1.25.8-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/exe /app/cmd/api

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/exe /app
COPY --from=builder /app/public /app/public
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]