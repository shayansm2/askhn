FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /api /worker ./
COPY internal/llm/prompts ./internal/llm/prompts
EXPOSE 8080
CMD ["/api"]
