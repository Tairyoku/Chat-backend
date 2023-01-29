FROM golang:latest
WORKDIR /app/src/chat-main
ENV GOPATH=/app
COPY ./ ./
#RUN #go get -u github.com/go-sql-driver/mysql
#RUN go get -u github.com/jinzhu/gorm
#RUN go get -u github.com/gorilla/mux
#RUN go get -u github.com/gorilla/handlers
RUN go mod download
RUN go build -o ./data ./cmd/main.go
CMD [ "./chat-main" ]