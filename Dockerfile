FROM node:24-alpine AS frontend-builder

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build-only

FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /static ./static

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /static ./static

EXPOSE 8080

CMD ["./main"]
