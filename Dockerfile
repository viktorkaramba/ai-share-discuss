FROM golang:1.21.5-alpine3.19 AS builder

COPY . /PlaylistsSynchronizer.Backend/
WORKDIR /PlaylistsSynchronizer.Backend/

#build app
RUN go mod download
RUN go build -o ./bin/api cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /PlaylistsSynchronizer.Backend/bin/api .
COPY --from=0 /PlaylistsSynchronizer.Backend/.env .
COPY --from=0 /PlaylistsSynchronizer.Backend/configs configs/

EXPOSE 8080

CMD ["./api"]