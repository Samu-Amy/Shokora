# Use Go
FROM golang:1.25.4-alpine

# Working Directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy other files
COPY *.go ./

# RUN go build -o /backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /backend


# Expose port 8080
EXPOSE 8080

# Run go executable file (from build)
CMD ["./backend"]

# TODO: fare la build fuori dal container (localmente o con GitHub Actions) e copiare nel container solo la build (i file generati) (usare crosscompiling e come funziona per i moduli terzi?)
