# Chat

A CRUD Golang API with MySQL & Docker

- Install
  * Go
  * Docker
  * docker-compose

- Usage: 

```bash
    cd $GOPATH/src

    git clone https://github.com/Tairyoku/Chat

    cd Chat

    docker-compose up
```

## Set up test database

```bash
    docker exec -it chat_db bash -l

    mysql -uroot -p@root
```

## Api endpoint out of container:

```bash
    cd $GOPATH/src/cmd/main.go
    
    go run main.go

```
- localhost:5050/api
