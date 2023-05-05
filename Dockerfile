# Use the official Golang image as a base image
FROM golang:1.17

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files, if they exist
COPY go.mod go.sum ./

# Set environment variables
ENV LDAP_BIND_USERNAME=UsernameFromGSuiteLdap \
    LDAP_BIND_PASSWORD=PasswordFromGSuiteLdap \
    LDAP_DC=dc=foo,dc=com \
    CRT_FILENAME=From_GSuite_LDAP.crt \
    KEY_FILENAME=From_GSuite_LDAP.key \
    LDAP_SERVER=ldap.google.com \
    LDAP_PORT=636 \
    RADIUS_SECRET=Long-Random-String-Probably-32-Characters \
    DEBUG=false

# Download dependencies
RUN go mod download

# Copy the rest of the application code, including the main.go file
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o gldapradius

# Start a new stage with a minimal image
FROM alpine:3.14

# Set the working directory
WORKDIR /app

# Copy the binary from the previous stage
COPY --from=0 /app/gldapradius /app/gldapradius

# Expose any required ports, if needed (optional)
EXPOSE 1812

# Run the application
CMD ["/app/gldapradius"]
