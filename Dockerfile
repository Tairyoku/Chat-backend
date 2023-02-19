FROM golang:latest
WORKDIR /app/server
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o docker-server ./cmd/main.go

EXPOSE 8080
CMD ["./docker-server"]
