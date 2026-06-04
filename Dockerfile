FROM node:24-alpine AS frontend-builder

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build-only

FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

COPY . .

COPY --from=frontend-builder /static ./static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=backend-builder /app/main .

COPY --from=frontend-builder /static ./static

EXPOSE 8080

CMD ["./main"]
