# Build stage
# Builder stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .

ARG user
ARG password
ARG host
ARG port
ARG dbname
ENV user="postgres"
ENV password="postgres"
ENV host="localhost"
ENV port="5432"
ENV dbname="quickkart"

#COPY /Users/rajatkr/GolandProjects/qvickly_new/pg.env /app
# ... rest of your dockerfile

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]