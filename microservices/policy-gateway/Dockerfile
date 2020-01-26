FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Remove the existing crt and key file
RUN rm server.crt
RUN rm server.key

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main .

# Command to run the executable
ENTRYPOINT ./main -listenPort ${LISTEN_PORT:-80} -forwardAddr ${FORWARD_ADDR} -forwardPort ${FORWARD_PORT} -policy ${POLICY_PATH}