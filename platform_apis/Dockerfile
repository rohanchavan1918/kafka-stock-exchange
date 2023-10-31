# Use the builder image
FROM golang:1.18 as builder
WORKDIR /app
COPY . .
RUN pwd
RUN ls
RUN ls ./config/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o go-boiler-binary

# Create a directory for the config volume
RUN mkdir /config

# multistage build
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/go-boiler-binary .

# Define a volume for the config directory
VOLUME ["/config"]

EXPOSE 8080

CMD [ "./go-boiler-binary" ]