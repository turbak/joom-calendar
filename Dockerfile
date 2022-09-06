FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/joom-calendar ./cmd/joom-calendar

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/bin/joom-calendar .

ENV PORT 8080
ENV DATABASE_URL "postgresql://postgres:postgrespw@localhost:49153/joom_calendar"
ENV GITHUB_CLIENT_ID "your_client_id"
ENV GITHUB_CLIENT_SECRET "your_client_secret"
ENV JWT_KEY "your_jwt_key"

EXPOSE $PORT

CMD ["./joom-calendar"]
