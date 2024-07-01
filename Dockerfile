FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache curl

ARG TARGETPLATFORM

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:3

ARG USERNAME=scorekeeper
ARG GROUP=scorekeeper

# Create the user
RUN addgroup -S $GROUP && adduser -S $USERNAME -G $GROUP

USER $USERNAME

WORKDIR /app
COPY --from=builder /app/main ./
COPY db/migration/ ./db/migration/

ENV GIN_MODE=release

EXPOSE 8080
CMD ["/app/main"]