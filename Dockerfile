# build stage
FROM golang:1.23.2-alpine AS builder

# install necessary build tools
RUN apk add --no-cache git make

# set working directory
WORKDIR /app

# install dependencies first (better cache utilization)
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build the application
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/agent main.go

# final/runtime stage
FROM scratch
#FROM alpine

# copy binary from builder
COPY --from=builder /app/agent /agent

# document exposed ports
EXPOSE 50001

# run the binary
CMD ["/agent"]