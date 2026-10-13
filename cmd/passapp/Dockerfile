FROM golang:1.27.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/passapp/exe /app/cmd/passapp

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/passapp/exe /app
CMD ["/app/exe"]