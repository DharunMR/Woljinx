# BUILD STAGE
From golang:1.25-alpine AS builder

WORKDIR /app

#Copy only mod file for caching
COPY go.mod go.sum ./
RUN go mod download

#Copy remaining files
COPY . .

RUN CGO-ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./server


# RUNTIME STAGE
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
