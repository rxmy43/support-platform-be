FROM golang:1.21

WORKDIR /app
COPY . /app

RUN apt-get update && apt-get install -y bash git curl && rm -rf /var/lib/apt/lists/*

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ENV PORT=8080
ENV DATABASE_URL=postgresql://neondb_owner:npg_soxp71WtwKNg@ep-old-frog-afa64s8n-pooler.c-2.us-west-2.aws.neon.tech/neondb?sslmode=require&channel_binding=require

RUN go build -o ./bin/server ./cmd/server/main.go

CMD migrate -path ./internal/db/migrations -database "$DATABASE_URL" up && ./bin/server

EXPOSE 8080
